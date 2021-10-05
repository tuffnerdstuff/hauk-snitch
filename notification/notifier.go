package notification

import (
	"fmt"
	"log"
	"net/smtp"
)

// Notifier can send email notifications about events in the mapper
type Notifier interface {
	NotifyNewSession(topic string, URL string)
	NotifyError(err interface{})
}

type notifier struct {
	config Config
}

// New returns a new Notifier instance
func New(config Config) Notifier {
	return &notifier{config: config}
}

func (t *notifier) NotifyNewSession(topic string, URL string) {
	t.sendMail(fmt.Sprintf("Forwarding %s to Hauk", topic), fmt.Sprintf("New session: %s", URL))
}

func (t *notifier) NotifyError(err interface{}) {
	t.sendMail("An error occurred", fmt.Sprintf("The following error occurred: %v", err))
}

func (t *notifier) sendMail(subject string, message string) {

	if t.config.Enabled {
		host := fmt.Sprintf("%s:%d", t.config.Host, t.config.Port)
		err := smtp.SendMail(host, nil, t.config.From, []string{t.config.To}, []byte(fmt.Sprintf("Subject: [hauk-snitch] %s\r\n\r\n%s", subject, message)))
		if err != nil {
			log.Printf("Could not send email notification: %v", err)
		}
	}
}
