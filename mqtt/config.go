package mqtt

const ParamType string = "_type"
const ParamLatitude string = "lat"
const ParamLongitude string = "lon"
const ParamAltitude string = "alt"
const ParamTime string = "tst"

// Config holds configuration for MqttClient
type Config struct {
	Host        string
	Port        int
	User        string
	Password    string
	IsTLS       bool
	IsAnonymous bool
}
