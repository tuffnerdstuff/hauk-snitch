package hauk

// Config for hauk backend
type Config struct {
	Host        string
	Port        int
	User        string
	Password    string
	IsTLS       bool
	IsAnonymous bool
	Duration    int
	Interval    int
}
