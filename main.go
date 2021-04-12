package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/tuffnerdstuff/hauk-snitch/config"
	"github.com/tuffnerdstuff/hauk-snitch/hauk"
	"github.com/tuffnerdstuff/hauk-snitch/mapper"
	"github.com/tuffnerdstuff/hauk-snitch/mqtt"
)

var mqttClient *mqtt.Client
var haukClient *hauk.Client

func main() {
	handleInterrupt()
	config.LoadConfig()

	runHaukClient()
	runMqttClient()

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
	log.Printf("%v\n", mqttConfig)
	mqttClient = mqtt.New(mqttConfig)
	mqttClient.Connect()

	for message := range mqttClient.Messages {
		location, err := mapper.CreateLocationFromMessage(message)
		if err != nil {
			log.Printf("Message does not contain valid location")
			continue
		}

		err = haukClient.PostLocation(message.Topic, location)
		if err != nil {
			log.Printf("Problem posting location: %v", err)
		}
	}
}

func runHaukClient() {
	haukConfig := config.GetHaukConfig()
	log.Printf("%v\n", haukConfig)
	haukClient = hauk.New(haukConfig)
}
