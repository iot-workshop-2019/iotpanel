package server

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	log "github.com/sirupsen/logrus"
)

// Server MQTT to manage device
type Server struct {
	Db     *sql.DB
	client mqtt.Client
}

var f mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("TOPIC: %s\n", msg.Topic())
	fmt.Printf("MSG: %s\n", msg.Payload())
}

// Publish to all device with MQTT
func (server *Server) Publish() {
	log.Info("qui..")
	if token := server.client.Publish("/radiolog/data", 0, false, "Server"); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}
}

// Init server module
func (server *Server) Init() {
	mqtt.ERROR = log.New()
	opts := mqtt.NewClientOptions().AddBroker("tcp://mqtt.asterix.cloud:1883").SetClientID("radiologHub")
	opts.SetKeepAlive(2 * time.Second)
	opts.SetDefaultPublishHandler(f)
	opts.SetPingTimeout(1 * time.Second)

	server.client = mqtt.NewClient(opts)
	if token := server.client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	if token := server.client.Subscribe("/radiolog/#", 0, nil); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}
}
