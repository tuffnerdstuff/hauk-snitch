package mapper

// NotificationConfig holds the configuration for the eMail notification about new Hauk sessions
type NotificationConfig struct {
	Enabled bool
	Host    string
	Port    int
	From    string
	To      string
}
