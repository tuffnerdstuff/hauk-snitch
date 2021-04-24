package notification

// Config holds the configuration for the eMail notification about new Hauk sessions
type Config struct {
	Enabled bool
	Host    string
	Port    int
	From    string
	To      string
}
