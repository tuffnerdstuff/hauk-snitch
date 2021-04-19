package mqtt

// Message is an MQTT message consisting of a topic and a body (map)
type Message struct {
	Topic string
	Body  map[string]interface{}
}
