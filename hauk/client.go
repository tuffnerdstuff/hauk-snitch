package hauk

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type client struct {
	config     Config
	httpClient http.Client
}

// Client is a client to the Hauk REST API
type Client interface {
	CreateSession() (Session, error)
	StopSession(sid string) error
	PostLocation(sid string, params url.Values) error
}

// New creates a new instance on a hauk client
func New(config Config) Client {
	httpClient := http.Client{}
	return &client{config: config, httpClient: httpClient}
}

// CreateSession attempts to create a new hauk session for a given device
func (t *client) CreateSession() (Session, error) {
	var session Session
	params := url.Values{
		"dur": {strconv.Itoa(t.config.Duration)},
		"int": {strconv.Itoa(t.config.Interval)},
		"pwd": {t.config.Password},
	}
	response, err := t.httpClient.PostForm(t.formatURL(EndpointCreate), params)

	err = getPostError(response, err, "posting session")
	if err != nil {
		return session, err
	}

	body, err := getBodyString(response)
	if err != nil {
		return session, err
	}

	bodyLines := strings.Split(body, "\n")
	session.ID = bodyLines[CreateResponseIndexID]
	session.SID = bodyLines[CreateResponseIndexSID]
	session.URL = bodyLines[CreateResponseIndexURL]

	return session, err

}

func (t *client) StopSession(sid string) error {

	// Set SID
	params := url.Values{}
	params.Add("sid", sid)

	// Send
	response, err := t.httpClient.PostForm(t.formatURL(EndpointStop), params)
	err = getPostError(response, err, "stopping session")
	if err != nil {
		return err
	}

	return err

}

// PostLocation sends a new location for the given device.
// If no previous session can be found, a new one will be created
func (t *client) PostLocation(sid string, params url.Values) error {

	// Add sid
	params.Add("sid", sid)

	// Send
	response, err := t.httpClient.PostForm(t.formatURL(EndpointPost), params)
	err = getPostError(response, err, "posting location")
	if err != nil {
		return err
	}
	// Parse body
	body, err := getBodyString(response)
	if err != nil {
		return err
	}

	// If session expired hauk returns Status 200 (OK) but "Session expired!"
	if strings.TrimSpace(body) == "Session expired!" {
		err = &SessionExpiredError{}
	}
	return err

}

func (t *client) formatURL(endpoint string) string {
	var protocol string
	if t.config.IsTLS {
		protocol = "https"
	} else {
		protocol = "http"
	}
	return fmt.Sprintf("%s://%s:%d/%s", protocol, t.config.Host, t.config.Port, endpoint)
}

func getBodyString(response *http.Response) (string, error) {
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		err = fmt.Errorf("Could not get session body: %w", err)
	}
	return string(body), err
}

func getPostError(response *http.Response, err error, action string) error {
	if err != nil {
		return fmt.Errorf("Error while %s: %w", action, err)
	} else if response.StatusCode != http.StatusOK {
		return fmt.Errorf("Server did not accept %s ( StatusCode = %d", action, response.StatusCode)
	}
	return nil
}
