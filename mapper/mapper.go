package mapper

import (
	"fmt"
	"log"
	"net/smtp"
	"net/url"
	"os"

	"github.com/mdp/qrterminal"
	"github.com/tuffnerdstuff/hauk-snitch/hauk"
	"github.com/tuffnerdstuff/hauk-snitch/mqtt"
)

// Mapper orchestrates incoming locations via mqtt and outgoing calls to Hauk
type Mapper struct {
	topicSessionMap    map[string]hauk.Session
	haukClient         hauk.Client
	notificationConfig NotificationConfig
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
func New(notificationConfig NotificationConfig, haukClient hauk.Client) Mapper {
	return Mapper{topicSessionMap: make(map[string]hauk.Session), haukClient: haukClient, notificationConfig: notificationConfig}
}

// Run maps mqtt messages to hauk API calls
func (t *Mapper) Run(messages <-chan mqtt.Message, haukClient hauk.Client) {

	for message := range messages {
		locationParams, err := createLocationParamsFromMessage(message)
		if err != nil {
			log.Printf("Message invalid, skipping: %s\n", err.Error())
			continue
		}

		var sid string
		if message.Body[mqtt.ParamTrigger] == mqtt.TriggerManual {
			sid, err = t.createNewSIDForTopic(message.Topic)
		} else {
			sid, err = t.getCurrentSIDForTopic(message.Topic)
		}
		if err != nil {
			log.Printf("%v\n", err.Error())
			continue
		}
		err = haukClient.PostLocation(sid, locationParams)
		if err != nil {
			switch err.(type) {
			case *hauk.SessionExpiredError:
				log.Printf("Session for %s expired, creating new one\n", message.Topic)
				var newSID string
				if newSID, err = t.createNewSIDForTopic(message.Topic); err != nil {
					log.Printf("%v", err.Error())
					continue
				}
				// re-send location
				log.Println("Re-posting location to new session")
				err = haukClient.PostLocation(newSID, locationParams)
				if err != nil {
					log.Printf("Re-posting failed, skipping: %v\n", err)
				}
			default:
				log.Printf("Invalid location %v: %v\n", locationParams, err.Error())
				continue
			}
		}
	}
}

func (t *Mapper) getCurrentSIDForTopic(topic string) (string, error) {
	session, sessionExists := t.topicSessionMap[topic]
	if !sessionExists {
		log.Printf("New topic %s, creating session\n", topic)
		return t.createNewSIDForTopic(topic)
	}
	return session.SID, nil
}

func (t *Mapper) createNewSIDForTopic(topic string) (string, error) {

	// Stop current session
	if currentSession, sessionExists := t.topicSessionMap[topic]; sessionExists {
		log.Printf("Stopping current session for %s: %v", topic, currentSession)
		err := t.haukClient.StopSession(currentSession.SID)
		if err != nil {
			log.Printf("Error while stopping current session %+v: %v", currentSession, err)
		}
	}

	// Create new Session
	newSession, err := t.haukClient.CreateSession()
	if err != nil {
		return "n/a", err
	}
	t.topicSessionMap[topic] = newSession

	// send email notification
	t.sendEmailNotification(topic, newSession.URL)

	// Print QR code on terminal
	log.Printf("New session for %s: %v", topic, newSession)
	qrterminal.GenerateHalfBlock(newSession.URL, qrterminal.L, os.Stdout)

	return newSession.SID, nil
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

func (t *Mapper) sendEmailNotification(topic string, URL string) {
	if t.notificationConfig.Enabled {
		host := fmt.Sprintf("%s:%d", t.notificationConfig.Host, t.notificationConfig.Port)
		err := smtp.SendMail(host, nil, t.notificationConfig.From, []string{t.notificationConfig.To}, []byte(fmt.Sprintf("Subject: Forwarding %s to Hauk\r\n\r\nNew session: %s", topic, URL)))
		if err != nil {
			log.Printf("Could not send email notification: %v", err)
		}
	}
}
