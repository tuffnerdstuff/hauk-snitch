package mqtt

const ParamType string = "_type"
const ParamLatitude string = "lat"
const ParamLongitude string = "lon"
const ParamAltitude string = "alt"
const ParamTime string = "tst"
const ParamAccuracy string = "acc"
const ParamVelocity string = "vel"

// Config holds configuration for MqttClient
type Config struct {
	Host        string
	Port        int
	User        string
	Password    string
	IsTLS       bool
	IsAnonymous bool
}
