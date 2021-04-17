package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/tuffnerdstuff/hauk-snitch/config"
	"github.com/tuffnerdstuff/hauk-snitch/hauk"
	m "github.com/tuffnerdstuff/hauk-snitch/mapper"
	"github.com/tuffnerdstuff/hauk-snitch/mqtt"
)

var mqttClient *mqtt.Client
var haukClient hauk.Client
var mapper m.Mapper

func main() {
	handleInterrupt()
	config.LoadConfig()

	runHaukClient()
	runMqttClient()
	runMapper()

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

func runMqttClient() {
	mqttConfig := config.GetMqttConfig()
	mqttClient = mqtt.New(mqttConfig)
	mqttClient.Connect()
}

func runHaukClient() {
	haukConfig := config.GetHaukConfig()
	haukClient = hauk.New(haukConfig)
}

func runMapper() {
	notificationConfig := config.GetNotificationConfig()
	mapper = m.New(notificationConfig, haukClient)
	mapper.Run(mqttClient.Messages, haukClient)
}
