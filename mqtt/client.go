package mqtt

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log"
	"time"

	paho "github.com/eclipse/paho.mqtt.golang"
)

// Client provides an mqtt client
type Client struct {
	Messages   chan Message
	config     Config
	pahoClient paho.Client
}

// New returns an instance of an mqtt client
func New(config Config) *Client {
	return &Client{config: config, Messages: make(chan Message)}
}

// Connect connects to mqtt broker using the given config
func (t *Client) Connect() {
	log.Printf("Connecting to mqtt broker %s\n", formatBrokerURL(t.config.Host, t.config.Port, t.config.IsTLS))
	t.initClient()
	t.connectClient()
	t.subscribeClient()
}

// Disconnect closes the Messages channel and disconnects the mqtt client
func (t *Client) Disconnect() {
	t.pahoClient.Disconnect(250)
	close(t.Messages)
}

func (t *Client) initClient() {
	opts := paho.NewClientOptions()
	opts.AddBroker(formatBrokerURL(t.config.Host, t.config.Port, t.config.IsTLS))
	opts.SetClientID("hauk-snitch" + generateHash())
	if !t.config.IsAnonymous {
		opts.SetUsername(t.config.User)
		opts.SetPassword(t.config.Password)
	}
	opts.SetCleanSession(false)
	// FIXME: process message
	opts.SetDefaultPublishHandler(func(client paho.Client, msg paho.Message) {
		jsonMap := make(map[string]interface{})
		json.Unmarshal(msg.Payload(), &jsonMap)
		t.Messages <- Message{Topic: msg.Topic(), Body: jsonMap}
	})
	t.pahoClient = paho.NewClient(opts)
}

func (t *Client) connectClient() {
	if token := t.pahoClient.Connect(); token.Wait() && token.Error() != nil {
		panic(fmt.Errorf("Error while connecting to mqtt broker: %w", token.Error()))
	}
}

func (t *Client) subscribeClient() {
	// FIXME: topic+qos configurable
	if token := t.pahoClient.Subscribe(t.config.Topic, byte(0), nil); token.Wait() && token.Error() != nil {
		panic(fmt.Errorf("Error while subscribing to topic: %w", token.Error()))
	}
}

func generateHash() string {
	sha := sha256.New()
	sha.Write([]byte(time.Now().String()))
	hashString := fmt.Sprintf("%x", sha.Sum(nil))
	return hashString[:10]
}

func formatBrokerURL(host string, port int, isTLS bool) string {
	var protocol string
	if isTLS {
		protocol = "ssl"
	} else {
		protocol = "tcp"
	}
	return fmt.Sprintf("%s://%s:%d", protocol, host, port)
}
