package server

import (
	"fmt"
	"time"

	"github.com/asterix24/radiolog-mqtt/dbi"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	log "github.com/sirupsen/logrus"
)

// Server MQTT to manage device
type Server struct {
	Db     *dbi.DBI
	client mqtt.Client
}

var f mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	log.Info(fmt.Sprintf("Recv: %s %s", msg.Topic(), msg.Payload()))
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
	opts.SetDefaultPublishHandler(f)
	opts.SetPingTimeout(1 * time.Second)

	server.client = mqtt.NewClient(opts)
	if token := server.client.Connect(); token.Wait() && token.Error() != nil {
		return token.Error()
	}

	if token := server.client.Subscribe("/radiolog/#", 0, nil); token.Wait() && token.Error() != nil {
		return token.Error()
	}

	return nil
}
