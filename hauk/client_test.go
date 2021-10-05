package hauk

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/mock"
)

type MockHTTPClient struct {
	mock.Mock
}

func (c *MockHTTPClient) PostForm(url string, data url.Values) (resp *http.Response, err error) {
	args := c.Called(url, data)
	return args.Get(0).(*http.Response), args.Error(1)
}

func TestCreateSession_HaukReturnsMalformedResponse_Error(t *testing.T) {

	// given: Parameters
	params := url.Values{
		"dur": {"47"},
		"int": {"11"},
		"pwd": {"pass"},
	}

	// given: HTTPClient which provides invalid return value for PostSession
	response := new(http.Response)
	response.StatusCode = http.StatusOK
	response.Body = ioutil.NopCloser(strings.NewReader("Something unexpected!"))
	httpClient := MockHTTPClient{}
	httpClient.On("PostForm", "http://www.example.local:4711/api/create.php", params).Return(response, nil).Once()

	// when
	client := New(Config{
		Host:        "www.example.local",
		Port:        4711,
		User:        "user",
		Password:    "pass",
		IsTLS:       false,
		IsAnonymous: true,
		Duration:    47,
		Interval:    11,
	}, &httpClient)
	_, err := client.CreateSession()

	// then: expect error
	_, ok := err.(*MalformedResponseError)
	if !ok {
		t.Fatal("Expected MalformedResponseError!")
	}

	// then: assert mock calls
	httpClient.AssertExpectations(t)
}
