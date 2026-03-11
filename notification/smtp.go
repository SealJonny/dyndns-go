package notification

import (
	"fmt"
	"log/slog"

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
	if client == nil {
		slog.Error("SetSMTPDefault called with nil SMTP client")
		return
	}
	defaultSMTPClient = *client
}

func EnableSMTP() {
	smtpEnabled = true
}

func (c *SMTPClient) sendEmail(subject, body string) {
	if !smtpEnabled {
		return
	}

	m := gomail.NewMessage()

	m.SetAddressHeader("From", c.user, "DynDNS Notifier")
	m.SetHeader("To", c.receiver)
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", body)

	d := gomail.NewDialer(c.host, c.port, c.user, c.password)
	if err := d.DialAndSend(m); err != nil {
		slog.Error("failed to send mail", "err", err, "subject", subject, "to", c.receiver)
	}
}

func (c *SMTPClient) Error(subject string, body string) {
	c.sendEmail(fmt.Sprintf("ERROR: %s", subject), body)
}

func SMTPError(subject string, body string) {
	defaultSMTPClient.Error(subject, body)
}

func (c *SMTPClient) Warn(subject string, body string) {
	c.sendEmail(fmt.Sprintf("WARN: %s", subject), body)
}

func SMTPWarn(subject string, body string) {
	defaultSMTPClient.Warn(subject, body)
}

func (c *SMTPClient) Info(subject string, body string) {
	c.sendEmail(fmt.Sprintf("INFO: %s", subject), body)
}

func SMTPInfo(subject string, body string) {
	defaultSMTPClient.Info(subject, body)
}
