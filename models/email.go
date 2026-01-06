package models

import (
	"fmt"

	"gopkg.in/gomail.v2"
)

const (
	DefaultSender = "support@lenslocked.com"
)

type Email struct {
	From      string
	To        string
	Subject   string
	Plaintext string
	HTML      string
}

type EmailService struct {
	DefaultSender string
	dialer        *gomail.Dialer
}

type SMTPConfig struct {
	Host     string
	Port     int
	Username string
	Password string
}

func NewEmailService(config SMTPConfig) *EmailService {
	es := EmailService{
		dialer: gomail.NewDialer(config.Host, config.Port, config.Username, config.Password),
	}
	return &es
}

func (es *EmailService) Send(email Email) error {
	msg := gomail.NewMessage()
	msg.SetHeader("To", email.To)
	// msg.SetHeader("From", email.From)
	es.setFrom(msg, email)
	msg.SetHeader("Subject", email.Subject)
	switch {
	case email.Plaintext != "" && email.HTML != "":
		msg.SetBody("test/plain", email.Plaintext)
		msg.AddAlternative("text/html", email.HTML)
	case email.Plaintext != "":
		msg.SetBody("test/plain", email.Plaintext)
	case email.HTML != "":
		msg.SetBody("test/html", email.HTML)
	}
	es.dialer.SSL = false
	err := es.dialer.DialAndSend(msg)
	if err != nil {
		return err
	}
	return nil
}

func (es *EmailService) setFrom(msg *gomail.Message, email Email) {
	var from string
	switch {
	case email.From != "":
		from = email.From
	case es.DefaultSender != "":
		from = es.DefaultSender
	default:
		from = DefaultSender
	}
	msg.SetHeader("From", from)
}

func (es *EmailService) ForgotPassword(to, resetURL string) error {
	email := Email{
		Subject:   "Reset your password",
		To:        to,
		Plaintext: "To reset your password, please visit the following link: " + resetURL,
		HTML:      `<p>Reset your password by click link here: <a href="` + resetURL + `">` + resetURL + `</a></p>`,
	}
	err := es.Send(email)
	if err != nil {
		return fmt.Errorf("forgot password email: %w", err)
	}
	return nil
}
