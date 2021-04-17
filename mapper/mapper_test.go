package mapper

import (
	"fmt"
	"net/url"
	"os"
	"reflect"
	"testing"

	"github.com/spf13/viper"
	"github.com/tuffnerdstuff/hauk-snitch/hauk"
	"github.com/tuffnerdstuff/hauk-snitch/mqtt"
)

type MockHaukClient struct {
	createSessionCalls int
	stopSessionCalls   int
	stopSessionSID     string
	postLocationCalls  int
	//postLocationSID    *string
	//postLocationParams *url.Values
}

func (t *MockHaukClient) CreateSession() (hauk.Session, error) {
	newSession := hauk.Session{ID: fmt.Sprintf("ID%d", t.createSessionCalls), SID: fmt.Sprintf("SID%d", t.createSessionCalls), URL: fmt.Sprintf("URL%d", t.createSessionCalls)}
	t.createSessionCalls++
	return newSession, nil
}

func (t *MockHaukClient) PostLocation(sid string, params url.Values) error {
	t.postLocationCalls++
	//t.postLocationSID = &sid
	//t.postLocationParams = &params
	return nil
}

func (t *MockHaukClient) StopSession(sid string) error {
	t.stopSessionSID = sid
	t.stopSessionCalls++
	return nil
}

func TestMain(m *testing.M) {

	// Disable email notification
	viper.SetDefault("notification.enabled", false)

	// Consume NewSessionsChannel
	go func() {
		for {
			<-NewSessionsChannel
		}
	}()

	os.Exit(m.Run())
}

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

func TestRun_LocationPushedAutomatically_KeepSession(t *testing.T) {

	// given: mqtt Location channel
	mqttLocations := make(chan mqtt.Message, 2)

	// given: First location update
	location := createValidLocationBody()
	mqttLocations <- mqtt.Message{Topic: "whatevs", Body: location}

	// given: Second (manual push) location update
	mqttLocations <- mqtt.Message{Topic: "whatevs", Body: location}

	close(mqttLocations)

	// given: Mock hauk client
	haukClient := MockHaukClient{stopSessionSID: "n/a"}

	// when
	Run(mqttLocations, &haukClient)

	// then
	expect(t, haukClient, MockHaukClient{
		createSessionCalls: 1,
		postLocationCalls:  2,
		stopSessionCalls:   0,
		stopSessionSID:     "n/a",
	})

}

func TestRun_LocationPushedManually_StopOldSessionAndStartNewSession(t *testing.T) {

	// given: mqtt Location channel
	mqttLocations := make(chan mqtt.Message, 2)

	// given: First location update
	location := createValidLocationBody()
	mqttLocations <- mqtt.Message{Topic: "whatevs", Body: location}

	// given: Second (manual push) location update
	location = createValidLocationBody()
	location["t"] = "u" // manual location push
	mqttLocations <- mqtt.Message{Topic: "whatevs", Body: location}

	close(mqttLocations)

	// given: Mock hauk client
	haukClient := MockHaukClient{}

	// when
	Run(mqttLocations, &haukClient)

	// then
	expect(t, haukClient, MockHaukClient{
		createSessionCalls: 2,
		postLocationCalls:  2,
		stopSessionCalls:   1,
		stopSessionSID:     "SID0",
	})

}

func expect(t *testing.T, actual interface{}, expected interface{}) {
	if actual != expected {
		t.Fatalf("actual value and expected value do not match:\nactual:%+v\nexpected:%+v", actual, expected)
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
