package notification

import (
	"fmt"

	"gopkg.in/gomail.v2"
)

var (
	defaultSMTPClient = SMTPClient{}
	smtpEnabled       = false
)

type SMTPClient struct {
	host     string
	port     int
	user     string
	password string
	receiver string
}

func New(host string, port int, user string, password string, receiver string) *SMTPClient {
	return &SMTPClient{
		host:     host,
		port:     port,
		user:     user,
		password: password,
		receiver: receiver,
	}
}

func SetSMTPDefault(client *SMTPClient) {
	defaultSMTPClient = *client
}

func EnableSMTP() {
	smtpEnabled = true
}

func (c *SMTPClient) sendEmail(subject, body string) error {
	if !smtpEnabled {
		return nil
	}

	m := gomail.NewMessage()

	m.SetHeader("From", c.user)
	m.SetHeader("To", c.receiver)
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", body)

	d := gomail.NewDialer(c.host, c.port, c.user, c.password)
	return d.DialAndSend(m)
}

func (c *SMTPClient) Error(subject string, body string) error {
	return c.sendEmail(fmt.Sprintf("ERROR: %s", subject), body)
}

func SMTPError(subject string, body string) error {
	return defaultSMTPClient.Error(subject, body)
}

func (c *SMTPClient) Warn(subject string, body string) error {
	return c.sendEmail(fmt.Sprintf("WARN: %s", subject), body)
}

func SMTPWarn(subject string, body string) error {
	return defaultSMTPClient.Warn(subject, body)
}

func (c *SMTPClient) Info(subject string, body string) error {
	return c.sendEmail(fmt.Sprintf("INFO: %s", subject), body)
}

func SMTPInfo(subject string, body string) error {
	return defaultSMTPClient.Info(subject, body)
}
