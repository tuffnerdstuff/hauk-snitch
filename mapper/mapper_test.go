package mapper

import (
	"fmt"
	"net/url"
	"reflect"
	"testing"

	"github.com/tuffnerdstuff/hauk-snitch/mqtt"
)

func TestMapMessageToLocation_InputOK_OutputOK(t *testing.T) {
	// given: mqtt message
	body := createValidLocationBody()
	// when
	givenLocation, _ := createLocationParamsFromMessage(mqtt.Message{Topic: "owntracks/user/device", Body: body})

	// then
	expectedLocation := url.Values{
		"lat":  {fmt.Sprintf("%v", body["lat"])},
		"lon":  {fmt.Sprintf("%v", body["lon"])},
		"acc":  {fmt.Sprintf("%v", body["acc"])},
		"alt":  {fmt.Sprintf("%v", body["alt"])},
		"spd":  {fmt.Sprintf("%v", body["vel"])},
		"time": {fmt.Sprintf("%v", int64(body["tst"].(float64)))},
	}
	if !reflect.DeepEqual(givenLocation, expectedLocation) {
		t.Fatalf("Locations do not match!\nGiven:%v\nExpected:%v", givenLocation, expectedLocation)
	}

}

func TestMapMessageToLocation_TypeNotLocation_Error(t *testing.T) {
	// given: type is not location
	body := make(map[string]interface{})
	body["_type"] = "somethingelse"

	// when
	_, err := createLocationParamsFromMessage(mqtt.Message{Topic: "owntracks/user/device", Body: body})

	// then: expect error
	if err == nil {
		t.Fatalf("Should return error")
	} else if err.Error() != "Type is not location" {
		t.Fatalf("Expected error informing user that _type != location")
	}

}

func createValidLocationBody() map[string]interface{} {
	body := make(map[string]interface{})
	body["_type"] = "location"
	body["acc"] = 5
	body["alt"] = 362
	body["batt"] = 76
	body["bs"] = 1
	body["conn"] = "w"
	body["created_at"] = 1.61826091e+09
	body["inregions"] = [...]string{"Home"}
	body["lat"] = 47.5968792
	body["lon"] = 12.9540961
	body["t"] = "p"
	body["tid"] = "dy"
	body["tst"] = 1.618243873e+09
	body["vac"] = 3
	body["vel"] = 42
	return body
}
