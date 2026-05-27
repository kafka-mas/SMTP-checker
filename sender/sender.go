package sender

import (
	"crypto/tls"
	"fmt"
	"os"

	"github.com/wneessen/go-mail"
)

type MailSender interface {
	//Send message to recipient
	Send(sendTo, subject, msg string) (string, error)
}

type mailSender struct {
	email      string
	password   string
	smtpServer string
	opts       options
}

type options struct {
	customAuth bool
	customPort int
	noTLS      bool
	helo       string
	username   string
}

type MailOptions func(*options)

/*
MailSender constructor

returns MailSender object
*/
func NewMailSender(login, password, smtpserver string, opts ...MailOptions) MailSender {
	helo, _ := os.Hostname()
	ms := &mailSender{
		email:      login,
		password:   password,
		smtpServer: smtpserver,
		opts: options{
			customPort: 465,
			noTLS:      false,
			helo:       helo,
			username:   login,
		},
	}

	for _, opt := range opts {
		opt(&ms.opts)
	}
	return ms
}

// Setup SMTP port
func WithCustomPort(port int) MailOptions {
	return func(o *options) {
		if port > 0 && port < 65535 {
			o.customPort = port
		}
	}
}

// Use an insecure connection without verifying the certificate
func WithNoTLS() MailOptions {
	return func(o *options) {
		o.noTLS = true
	}
}

// Setup HELO
func WithHELO(helo string) MailOptions {
	return func(o *options) {
		o.helo = helo
	}
}

// set username values if it is not specified in the email@domain format
func WithCustomUsername(username string) MailOptions {
	return func(o *options) {
		o.username = username
	}
}

func (m *mailSender) Send(sendTo, subject, msg string) (string, error) {
	message := mail.NewMsg()
	if err := message.From(m.email); err != nil {
		return "", fmt.Errorf("failed to set From address %v: %w", m.email, err)
	}
	if err := message.To(sendTo); err != nil {
		return "", fmt.Errorf("failed to set To address %v: %w", sendTo, err)
	}
	message.Subject(subject)
	message.SetBodyString(mail.TypeTextPlain, msg)

	client, err := func(isNotSecure bool) (*mail.Client, error) {
		if isNotSecure {
			return mail.NewClient(m.smtpServer,
				mail.WithHELO(m.opts.helo),
				mail.WithUsername(m.opts.username),
				mail.WithPassword(m.password),
				mail.WithSMTPAuth(mail.SMTPAuthLogin),
				mail.WithTLSConfig(&tls.Config{InsecureSkipVerify: true}),
				mail.WithPort(m.opts.customPort),
			)
		} else {
			return mail.NewClient(m.smtpServer,
				mail.WithSMTPAuth(mail.SMTPAuthAutoDiscover),
				mail.WithUsername(m.opts.username),
				mail.WithPassword(m.password),
				mail.WithHELO(m.opts.helo),
			)
		}
	}(m.opts.noTLS)
	if err != nil {
		return "", fmt.Errorf("failed to create mail client: %w", err)
	}
	defer client.Close()

	if err := client.DialAndSend(message); err != nil {
		return "", fmt.Errorf("failed to send mail: %w", err)
	}

	messageID := message.GetGenHeader(mail.HeaderMessageID)[0]

	return messageID, nil
}
