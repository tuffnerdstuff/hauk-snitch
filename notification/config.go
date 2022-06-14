package notification

// Config holds the configuration for the eMail notification about new Hauk sessions
type Config struct {
	Smtp   SMTPConfig
	Gotify GotifyConfig
}

type SMTPConfig struct {
	Enabled bool
	Host    string
	Port    int
	From    string
	To      string
}

type GotifyConfig struct {
	Enabled  bool
	URL      string
	AppToken string
	Priority int
}
