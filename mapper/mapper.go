package mapper

import (
	"fmt"
	"log"
	"net/url"

	"github.com/tuffnerdstuff/hauk-snitch/hauk"
	"github.com/tuffnerdstuff/hauk-snitch/mqtt"
	"github.com/tuffnerdstuff/hauk-snitch/notification"
)

// Mapper orchestrates incoming locations via mqtt and outgoing calls to Hauk
type Mapper struct {
	haukClient hauk.Client
	notifier   notification.Notifier
	config     Config
}

type valueMapping struct {
	haukKey         string
	mappingFunction func(value interface{}) string
}

var mqttToHaukKeyMap = map[string]valueMapping{
	mqtt.ParamAccuracy:  {hauk.ParamAccuracy, nil},
	mqtt.ParamAltitude:  {hauk.ParamAltitude, nil},
	mqtt.ParamLatitude:  {hauk.ParamLatitude, nil},
	mqtt.ParamLongitude: {hauk.ParamLongitude, nil},
	mqtt.ParamVelocity: {hauk.ParamVelocity, func(value interface{}) string {
		// km/h -> m/s
		return fmt.Sprintf("%f", convertToFloat(value)/3.6)
	}},
	mqtt.ParamTime: {hauk.ParamTime, func(value interface{}) string {
		// UNIX epoch float -> int
		// Hauk Android client also sends float, but formatted differently.
		// Before converting to int the frontend sometimes did not update.
		return fmt.Sprintf("%d", int64(convertToFloat(value)))
	}},
}

// New creates a new instance of the mapper orchestrating mqtt and Hauk
func New(config Config, haukClient hauk.Client, notifier notification.Notifier) Mapper {
	return Mapper{haukClient: haukClient, config: config, notifier: notifier}
}

// Run maps mqtt messages to hauk API calls
func (t *Mapper) Run(messages <-chan mqtt.Message) {

	topicChannels := make(map[string]chan mqtt.Message)
	for message := range messages {

		topicChannel, available := topicChannels[message.Topic]
		if !available {
			topicChannel = make(chan mqtt.Message, 1)
			topicChannels[message.Topic] = topicChannel
		}
		topicChannel <- message
		go t.handleTopic(topicChannel)

	}
}

func (t *Mapper) handleTopic(topicMessages <-chan mqtt.Message) {
	sid := ""
	for message := range topicMessages {
		locationParams, err := createLocationParamsFromMessage(message)
		if err != nil {
			log.Printf("Message invalid, skipping: %s\n", err.Error())
			continue
		}

		for {
			var newSID string = sid
			if sid == "" && t.config.SessionStartAuto {
				log.Printf("New topic %s, creating session\n", message.Topic)
				newSID, err = t.createNewSIDForTopic(message.Topic)
			} else if t.config.SessionStartManual && message.Body[mqtt.ParamTrigger] == mqtt.TriggerManual {
				log.Printf("Manual location push, creating new session for topic %s\n", message.Topic)
				t.stopSession(sid)
				newSID, err = t.createNewSIDForTopic(message.Topic)
			}
			if err != nil {
				log.Printf("%v\n", err.Error())
				break
			} else if newSID == "" {
				log.Printf("Starting a session not allowed for location: %v\n", message)
				break
			}
			sid = newSID

			err = t.haukClient.PostLocation(sid, locationParams)
			if err == nil {
				break
			} else if !isSessionExpired(err) {
				log.Printf("Error while sending location, discarding location: %s\n", err.Error())
				break
			}
			// session expired, resetting sid
			sid = ""
		}
	}

}

func (t *Mapper) stopSession(sid string) {
	// Stop current session
	if t.config.SessionStopAuto {
		if sid != "" {
			log.Printf("Stopping current session %s\n", sid)
			err := t.haukClient.StopSession(sid)
			if err != nil {
				log.Printf("Error while stopping current session %+v: %v\n", sid, err)
			}
		}
	}
}

func (t *Mapper) createNewSIDForTopic(topic string) (string, error) {

	// Create new Session
	newSession, err := t.haukClient.CreateSession()
	if err != nil {
		return "n/a", err
	}

	// send email notification
	t.notifier.NotifyNewSession(topic, newSession.URL)

	log.Printf("New session for %s: %v\n", topic, newSession)

	return newSession.SID, nil
}

func isSessionExpired(err error) bool {
	if err != nil {
		switch err.(type) {
		case *hauk.SessionExpiredError:
			return true
		}
	}
	return false
}

func createLocationParamsFromMessage(msg mqtt.Message) (url.Values, error) {
	body := msg.Body
	haukValues := url.Values{}
	if body[mqtt.ParamType] == "location" {
		for mqttKey, mqttValue := range body {
			setHaukValue(&haukValues, mqttKey, mqttValue)

		}
		return haukValues, nil
	}
	return haukValues, fmt.Errorf("Type is not location")

}

func setHaukValue(haukValues *url.Values, key string, value interface{}) {
	valueMapping, hasMapping := mqttToHaukKeyMap[key]
	if hasMapping {
		var convertedValue string
		if valueMapping.mappingFunction != nil {
			convertedValue = valueMapping.mappingFunction(value)
		} else {
			convertedValue = fmt.Sprintf("%v", value)
		}
		haukValues.Add(valueMapping.haukKey, convertedValue)
	}
}

func convertToFloat(value interface{}) float64 {
	var floatValue float64 = 0
	switch value.(type) {
	case float64:
		floatValue = value.(float64)
	case int:
		floatValue = float64(value.(int))
	case int32:
		floatValue = float64(value.(int32))
	case int64:
		floatValue = float64(value.(int64))
	}
	return floatValue
}
