package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime/debug"

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
	defer handlePanic()
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

func handlePanic() {
	if err := recover(); err != nil {
		message := fmt.Sprintf("hauk-snitch panicked and terminated!\nerror: %v\nstacktrace: %s\n", err, string(debug.Stack()))
		log.Print(message)
		if notifier != nil {
			notifier.NotifyError(message)
		}
	}
}

func initMqttClient() {
	mqttClient = mqtt.New(config.GetMqttConfig())
	mqttClient.Connect()
}

func initHaukClient() {
	haukClient = hauk.New(config.GetHaukConfig(), &http.Client{})
}

func initNotifier() {
	notifier = notification.New(config.GetNotificationConfig())
}

func initMapper() {
	mapper = m.New(config.GetMapperConfig(), haukClient, notifier)
	mapper.Run(mqttClient.Messages)
}
