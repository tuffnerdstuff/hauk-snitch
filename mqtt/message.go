package mqtt

type Message struct {
	Topic string
	Body  map[string]interface{}
}
