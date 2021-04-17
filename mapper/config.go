package mapper

type NotificationConfig struct {
	Enabled bool
	Host    string
	Port    int
	From    string
	To      string
}
