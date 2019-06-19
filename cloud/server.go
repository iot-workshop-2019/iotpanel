package server

import (
	"fmt"
	"os"
	"regexp"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/iot-workshop-2019/iotpanel/dbi"
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
	Db     *dbi.DBI
	client mqtt.Client
}

func (server *Server) f(client mqtt.Client, msg mqtt.Message) {
	log.Info(fmt.Sprintf("Recv: %s %s", msg.Topic(), msg.Payload()))
	re := regexp.MustCompile(`(Node.[a-zA-Z0-9]{6})/(.*)$`)
	topic := re.FindStringSubmatch(msg.Topic())

	log.Infof("Topics[%v]: payload[%s]", topic, msg.Payload())
	if len(topic) < 2 {
		return
	}
	server.Db.UpdateNode(topic[1], topic[2], string(msg.Payload()))
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
	opts := mqtt.NewClientOptions()
	nmclientid, ok := os.LookupEnv("NMCLIENTID")
	if !ok {
		log.Error("NMCLIENTID environment variable required but not set")
		log.Info("Using default variable radiologhub")
		nmclientid = "radiologhub"
	}
	opts.AddBroker("tcp://mqtt.asterix.cloud:1883").SetClientID(nmclientid)
	opts.SetKeepAlive(time.Second * time.Duration(60))
	opts.SetConnectionLostHandler(func(client mqtt.Client, e error) {
		log.Warn(fmt.Sprintf("Connection lost : %v", e))
		if client.IsConnected() {
			client.Disconnect(500)
		}
		server.client = mqtt.NewClient(opts)
		connect(server.client)
		subscribe(server.client, server.f)
	})

	server.client = mqtt.NewClient(opts)
	connect(server.client)
	subscribe(server.client, server.f)
	log.Info("Start to subscribe....")
	defer func() {
		if r := recover(); r != nil {
			log.Warn("unknown panic error, try to recover connection,", r)
			connect(server.client)
			subscribe(server.client, server.f)
		}
	}()

	return nil
}

func connect(client mqtt.Client) {
	for {
		// create connection
		token := client.Connect()
		if token.Wait() && token.Error() != nil {
			log.Errorf("Fail to connect broker, %v", token.Error())
			time.Sleep(5 * time.Second)

			log.Error("Retry the connection")
			continue
		} else {
			log.Error("Reconnection successful!")
			break
		}
	}
}

func subscribe(client mqtt.Client, onIncomingDataReceived mqtt.MessageHandler) {
	for {
		// subscribe the topic, "#" means all topics
		token := client.Subscribe("#", byte(0), onIncomingDataReceived)
		if token.Wait() && token.Error() != nil {
			log.Error("Fail to sub... ", token.Error())
			time.Sleep(5 * time.Second)

			log.Errorf("Retry to subscribe")
			continue
		} else {
			log.Info("Subscribe successful!")
			break
		}
	}
}
