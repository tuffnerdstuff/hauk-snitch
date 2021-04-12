package mapper

import (
	"fmt"

	"github.com/tuffnerdstuff/hauk-snitch/hauk"
	"github.com/tuffnerdstuff/hauk-snitch/mqtt"
)

func CreateLocationFromMessage(msg mqtt.Message) (hauk.Location, error) {
	body := msg.Body
	if body[mqtt.ParamType] == "location" {
		return hauk.Location{
			Latitude:  body[mqtt.ParamLatitude].(float64),
			Longitude: body[mqtt.ParamLongitude].(float64),
			Time:      body[mqtt.ParamTime].(float64),
		}, nil
	}
	return hauk.Location{}, fmt.Errorf("Type is not location")

}
