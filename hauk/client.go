package hauk

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// Client is a hauk client
type Client struct {
	config     Config
	httpClient http.Client
}

// New creates a new instance on a hauk client
func New(config Config) *Client {
	httpClient := http.Client{}
	return &Client{config: config, httpClient: httpClient}
}

// CreateSession attempts to create a new hauk session for a given device
func (t *Client) CreateSession() (Session, error) {
	var session Session
	params := url.Values{
		"dur": {strconv.Itoa(t.config.Duration)},
		"int": {strconv.Itoa(t.config.Interval)},
		"pwd": {t.config.Password},
	}
	response, err := t.httpClient.PostForm(t.formatURL(EndpointCreate), params)

	err = getPostError(response, err, "session")
	if err != nil {
		return session, err
	}

	body, err := getBodyString(response)
	if err != nil {
		return session, fmt.Errorf("Could not get session body")
	}

	bodyLines := strings.Split(body, "\n")
	session.ID = bodyLines[CreateResponseIndexID]
	session.SID = bodyLines[CreateResponseIndexSID]
	session.URL = bodyLines[CreateResponseIndexURL]

	return session, err

}

// PostLocation sends a new location for the given device.
// If no previous session can be found, a new one will be created
func (t *Client) PostLocation(sid string, params url.Values) error {

	// Add sid
	params.Add("sid", sid)

	// Send
	response, err := t.httpClient.PostForm(t.formatURL(EndpointPost), params)
	err = getPostError(response, err, "location")
	return err

}

func (t *Client) formatURL(endpoint string) string {
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
	return string(body), err
}

func getPostError(response *http.Response, err error, entity string) error {
	if err != nil {
		return fmt.Errorf("Error while posting %s: %w", entity, err)
	} else if response.StatusCode != http.StatusOK {
		return fmt.Errorf("Server did not accept %s ( StatusCode = %d", entity, response.StatusCode)
	}
	return nil
}
