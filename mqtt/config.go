package mqtt

// Config holds configuration for MqttClient
type Config struct {
	Host        string
	Port        int
	User        string
	Password    string
	IsTLS       bool
	IsAnonymous bool
}
