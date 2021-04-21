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
	mqttLocations := make(chan mqtt.Message, 1)
	mqttLocations <- mqtt.Message{Topic: "whatevs", Body: body}
	close(mqttLocations)

	// given: Mock hauk client
	haukClient := new(MockHaukClient)

	// given: notifier
	notifier := new(MockNotifier)

	// when
	mapper := New(Config{
		SessionStartAuto:   true,
		SessionStartManual: true,
		SessionStopAuto:    true,
	}, haukClient, notifier)
	mapper.Run(mqttLocations)

	// then: assert mock calls
	haukClient.AssertExpectations(t)
	notifier.AssertExpectations(t)

}

func TestRun_SessionAutoStartAndManualStartAndAutoStop(t *testing.T) {
	testSessionHandling(t, true, true, true)
}

func TestRun_NoSessionManualStart(t *testing.T) {
	testSessionHandling(t, true, false, true)
}

func TestRun_NoSessionAutoStart(t *testing.T) {
	testSessionHandling(t, false, true, true)
}

func TestRun_SessionManualStartOnly(t *testing.T) {
	testSessionHandling(t, false, true, false)
}

func TestRun_NoSessionAutoStop(t *testing.T) {
	testSessionHandling(t, true, true, false)
}

func TestRun_SessionAutoStartOnly(t *testing.T) {
	testSessionHandling(t, true, false, false)
}

func TestRun_NoSession(t *testing.T) {
	testSessionHandling(t, false, false, false)
	testSessionHandling(t, false, false, true)
}

func testSessionHandling(t *testing.T, startSessionAuto bool, startSessionManual bool, stopSessionAuto bool) {

	// given: valid locations
	locationAuto := createValidLocationBody()
	locationAuto["tst"] = float64(1)
	locationManual := createValidLocationBody()
	locationManual["t"] = "u"
	locationManual["tst"] = float64(2)

	// given: Mock hauk client
	haukClient := new(MockHaukClient)

	// given: notifier
	notifier := new(MockNotifier)

	// auto push locationAuto
	if startSessionAuto {
		// --> CreateSession "firstSession"
		haukClient.On("CreateSession").Return(hauk.Session{SID: "firstSession", URL: "firstURL"}, nil).Once()
		notifier.On("NotifyNewSession", "whatevs", "firstURL").Once()
		// --> PostLocation to "firstSession"
		haukClient.On("PostLocation", "firstSession", getExpectedLocationValues(locationAuto)).Return(&hauk.SessionExpiredError{}).Once()
		// handle expired session
		// --> CreateSession "secondSession"
		haukClient.On("CreateSession").Return(hauk.Session{SID: "secondSession", URL: "secondURL"}, nil).Once()
		notifier.On("NotifyNewSession", "whatevs", "secondURL").Once()
		// --> PostLocation to "secondSession" (re-send)
		haukClient.On("PostLocation", "secondSession", getExpectedLocationValues(locationAuto)).Return(nil).Once()
	}
	// manual push locationManual
	if startSessionManual {
		if startSessionAuto && stopSessionAuto {
			// --> StopSession "secondSession"
			haukClient.On("StopSession", "secondSession").Return(nil).Once()
		}
		// --> CreateSession "thirdSession"
		haukClient.On("CreateSession").Return(hauk.Session{SID: "thirdSession", URL: "thirdURL"}, nil).Once()
		notifier.On("NotifyNewSession", "whatevs", "thirdURL").Once()
		// --> PostLocation to "thirdSession"
		haukClient.On("PostLocation", "thirdSession", getExpectedLocationValues(locationManual)).Return(nil).Once()
	} else if startSessionAuto {
		// --> PostLocation to "secondSession"
		haukClient.On("PostLocation", "secondSession", getExpectedLocationValues(locationManual)).Return(nil).Once()
	}

	// given: mqtt Location channel
	mqttLocations := make(chan mqtt.Message, 2)

	// given: First location update
	mqttLocations <- mqtt.Message{Topic: "whatevs", Body: locationAuto}

	// given: Second location update
	mqttLocations <- mqtt.Message{Topic: "whatevs", Body: locationManual}

	close(mqttLocations)

	// when
	mapper := New(Config{
		SessionStartAuto:   startSessionAuto,
		SessionStartManual: startSessionManual,
		SessionStopAuto:    stopSessionAuto,
	}, haukClient, notifier)
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
