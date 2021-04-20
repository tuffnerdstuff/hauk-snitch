package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/tuffnerdstuff/hauk-snitch/config"
	"github.com/tuffnerdstuff/hauk-snitch/hauk"
	m "github.com/tuffnerdstuff/hauk-snitch/mapper"
	"github.com/tuffnerdstuff/hauk-snitch/mqtt"
	"github.com/tuffnerdstuff/hauk-snitch/notification"
)

var mqttClient *mqtt.Client
var haukClient hauk.Client
var notifier notification.Notifier
var mapper m.Mapper

func main() {
	handleInterrupt()
	config.LoadConfig()

	initHaukClient()
	initMqttClient()
	initNotifier()
	initMapper()

}

func handleInterrupt() {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, os.Kill)
	go func() {
		<-interrupt
		log.Println("Exiting")
		if mqttClient != nil {
			mqttClient.Disconnect()
		}
	}()
}

func initMqttClient() {
	mqttClient = mqtt.New(config.GetMqttConfig())
	mqttClient.Connect()
}

func initHaukClient() {
	haukClient = hauk.New(config.GetHaukConfig())
}

func initNotifier() {
	notifier = notification.New(config.GetNotificationConfig())
}

func initMapper() {
	mapper = m.New(config.GetMapperConfig(), haukClient, notifier)
	mapper.Run(mqttClient.Messages)
}
