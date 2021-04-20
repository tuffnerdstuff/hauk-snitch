package mapper

import (
	"fmt"
	"net/url"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/tuffnerdstuff/hauk-snitch/hauk"
	"github.com/tuffnerdstuff/hauk-snitch/mqtt"
)

type MockHaukClient struct {
	mock.Mock
}

func (t *MockHaukClient) CreateSession() (hauk.Session, error) {
	args := t.Called()
	return args.Get(0).(hauk.Session), args.Error(1)
}

func (t *MockHaukClient) PostLocation(sid string, params url.Values) error {
	args := t.Called(sid, params)
	return args.Error(0)
}

func (t *MockHaukClient) StopSession(sid string) error {
	args := t.Called(sid)
	return args.Error(0)
}

type MockNotifier struct {
	mock.Mock
}

func (t *MockNotifier) NotifyNewSession(topic string, URL string) {
	t.Called(topic, URL)
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

	// given: valid location
	location := createValidLocationBody()

	// given: Mock hauk client
	haukClient := new(MockHaukClient)
	haukClient.On("CreateSession").Return(hauk.Session{SID: "newSessionSID", URL: "url/to/newSessionSID"}, nil)
	haukClient.On("PostLocation", "newSessionSID", getExpectedLocationValues(location)).Return(nil)

	// given: mqtt Location channel
	mqttLocations := make(chan mqtt.Message, 2)

	// given: First location update
	mqttLocations <- mqtt.Message{Topic: "whatevs", Body: location}

	// given: Second location update
	mqttLocations <- mqtt.Message{Topic: "whatevs", Body: location}

	close(mqttLocations)

	// given: notifier
	notifier := new(MockNotifier)
	notifier.On("NotifyNewSession", "whatevs", "url/to/newSessionSID").Once()

	// when
	mapper := New(Config{}, haukClient, notifier)
	mapper.Run(mqttLocations)

	// then: Assert haukClient Calls
	haukClient.AssertExpectations(t)
	haukClient.AssertNumberOfCalls(t, "CreateSession", 1)
	haukClient.AssertNumberOfCalls(t, "PostLocation", 2)
	haukClient.AssertNotCalled(t, "StopSession", "newSessionID")

	// then: Assert notifier Calls
	notifier.AssertExpectations(t)

}

func TestRun_LocationPushedAutomaticallyButSessionExpired_StartNewSession(t *testing.T) {

	// given: valid locations
	location1 := createValidLocationBody()
	location1["tst"] = float64(1)
	location2 := createValidLocationBody()
	location2["tst"] = float64(2)

	// given: Mock hauk client
	haukClient := new(MockHaukClient)
	haukClient.On("CreateSession").Return(hauk.Session{SID: "firstSession", URL: "firstURL"}, nil).Once()
	haukClient.On("PostLocation", "firstSession", getExpectedLocationValues(location1)).Return(&hauk.SessionExpiredError{}).Once()
	haukClient.On("CreateSession").Return(hauk.Session{SID: "secondSession", URL: "secondURL"}, nil).Once()
	haukClient.On("PostLocation", "secondSession", getExpectedLocationValues(location1)).Return(nil).Once()
	haukClient.On("PostLocation", "secondSession", getExpectedLocationValues(location2)).Return(nil).Once()

	// given: mqtt Location channel
	mqttLocations := make(chan mqtt.Message, 2)

	// given: First location update
	mqttLocations <- mqtt.Message{Topic: "whatevs", Body: location1}

	// given: Second location update
	mqttLocations <- mqtt.Message{Topic: "whatevs", Body: location2}

	close(mqttLocations)

	// given: notifier
	notifier := new(MockNotifier)
	notifier.On("NotifyNewSession", "whatevs", "firstURL").Once()
	notifier.On("NotifyNewSession", "whatevs", "secondURL").Once()

	// when
	mapper := New(Config{}, haukClient, notifier)
	mapper.Run(mqttLocations)

	// then: assert mock calls
	haukClient.AssertExpectations(t)
	notifier.AssertExpectations(t)

}

func TestRun_LocationPushedManually_StopOldSessionAndStartNewSession(t *testing.T) {

	// given: valid locations
	locationAuto := createValidLocationBody()
	locationAuto["tst"] = float64(1)
	locationManual := createValidLocationBody()
	locationManual["tst"] = float64(2)
	locationManual["t"] = "u"

	// given: Mock hauk client
	haukClient := new(MockHaukClient)
	haukClient.On("CreateSession").Return(hauk.Session{SID: "firstSession", URL: "firstURL"}, nil).Once()
	haukClient.On("PostLocation", "firstSession", getExpectedLocationValues(locationAuto)).Return(nil).Once()
	haukClient.On("StopSession", "firstSession").Return(nil).Once()
	haukClient.On("CreateSession").Return(hauk.Session{SID: "secondSession", URL: "secondURL"}, nil).Once()
	haukClient.On("PostLocation", "secondSession", getExpectedLocationValues(locationManual)).Return(nil).Once()

	// given: mqtt Location channel
	mqttLocations := make(chan mqtt.Message, 2)

	// given: First location update
	mqttLocations <- mqtt.Message{Topic: "whatevs", Body: locationAuto}

	// given: Second location update
	mqttLocations <- mqtt.Message{Topic: "whatevs", Body: locationManual}

	close(mqttLocations)

	// given: notifier
	notifier := new(MockNotifier)
	notifier.On("NotifyNewSession", "whatevs", "firstURL").Once()
	notifier.On("NotifyNewSession", "whatevs", "secondURL").Once()

	// when
	mapper := New(Config{}, haukClient, notifier)
	mapper.Run(mqttLocations)

	// then: assert mock calls
	haukClient.AssertExpectations(t)
	notifier.AssertExpectations(t)

}

func getExpectedLocationValues(location map[string]interface{}) url.Values {
	return url.Values{
		"lat":  {fmt.Sprintf("%v", location["lat"])},
		"lon":  {fmt.Sprintf("%v", location["lon"])},
		"acc":  {fmt.Sprintf("%v", location["acc"])},
		"alt":  {fmt.Sprintf("%v", location["alt"])},
		"spd":  {fmt.Sprintf("%f", float64(location["vel"].(int))/3.6)},
		"time": {fmt.Sprintf("%d", int64(location["tst"].(float64)))},
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
