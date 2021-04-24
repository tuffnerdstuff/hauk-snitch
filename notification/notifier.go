package notification

import (
	"fmt"
	"log"
	"net/smtp"
)

// Notifier can send email notifications about events in the mapper
type Notifier interface {
	NotifyNewSession(topic string, URL string)
}

type notifier struct {
	config Config
}

// New returns a new Notifier instance
func New(config Config) Notifier {
	return &notifier{config: config}
}

func (t *notifier) NotifyNewSession(topic string, URL string) {
	if t.config.Enabled {
		host := fmt.Sprintf("%s:%d", t.config.Host, t.config.Port)
		err := smtp.SendMail(host, nil, t.config.From, []string{t.config.To}, []byte(fmt.Sprintf("Subject: Forwarding %s to Hauk\r\n\r\nNew session: %s", topic, URL)))
		if err != nil {
			log.Printf("Could not send email notification: %v", err)
		}
	}
}
