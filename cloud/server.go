package server

import (
	"fmt"
	"regexp"
	"time"

	"github.com/asterix24/radiolog-mqtt/dbi"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	log "github.com/sirupsen/logrus"
)

// MsgFmt ...
type MsgFmt struct {
	Timestamp string `json:"timestamp"`
	Node      string `json:"node"`
	Data      string `json:"data"`
}

// Server MQTT to manage device
type Server struct {
	Db       *dbi.DBI
	StatusEv chan MsgFmt
	DataEv   chan MsgFmt
	client   mqtt.Client
}

var statusEv = make(chan MsgFmt)
var dataEv = make(chan MsgFmt)

func (server *Server) f(client mqtt.Client, msg mqtt.Message) {
	log.Info(fmt.Sprintf("Recv: %s %s", msg.Topic(), msg.Payload()))
	re := regexp.MustCompile(`(Node-[a-zA-Z0-9]{6})/(status)$`)
	topic := re.FindStringSubmatch(msg.Topic())

	if len(topic) < 2 {
		return
	}

	what := topic[2]
	switch what {
	case "status":
		statusEv <- MsgFmt{Node: topic[1], Data: string(msg.Payload()), Timestamp: time.Now().Format("2006-01-02 15:04:05")}
	case "data":
		dataEv <- MsgFmt{Node: topic[1], Data: string(msg.Payload()), Timestamp: time.Now().Format("2006-01-02 15:04:05")}
	default:
		log.Infof("%s: %s", topic[1], msg.Payload())
	}
}

// Publish to all device with MQTT
func (server *Server) Publish(key string, value string) error {
	data := fmt.Sprintf("/radiolog/%s", key)
	if token := server.client.Publish(data, 0, false, value); token.Wait() && token.Error() != nil {
		return token.Error()
	}

	return nil
}

// Init server module
func (server *Server) Init() error {
	mqtt.ERROR = log.New()
	opts := mqtt.NewClientOptions().AddBroker("tcp://mqtt.asterix.cloud:1883").SetClientID("radiologHub")
	opts.SetKeepAlive(2 * time.Second)
	opts.SetDefaultPublishHandler(server.f)
	opts.SetPingTimeout(1 * time.Second)

	server.client = mqtt.NewClient(opts)
	if token := server.client.Connect(); token.Wait() && token.Error() != nil {
		return token.Error()
	}

	if token := server.client.Subscribe("/radiolog/#", 0, nil); token.Wait() && token.Error() != nil {
		return token.Error()
	}

	server.StatusEv = statusEv
	server.DataEv = dataEv

	return nil
}
