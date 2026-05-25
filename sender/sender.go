package sender

import (
	"crypto/tls"
	"log"
	"os"

	"github.com/wneessen/go-mail"
)

type MailSender interface {
	Send(sendTo, subject, msg string) ([]string, error)
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

func NewMailSender(login, password, smtpserver string, opts ...MailOptions) MailSender {
	helo, _ := os.Hostname()
	ms := &mailSender{
		email:      login,
		password:   password,
		smtpServer: smtpserver,
		opts: options{
			customAuth: false,
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

func WithCustomAuth() MailOptions {
	return func(o *options) {
		o.customAuth = true
	}
}

func WithCustomPort(port int) MailOptions {
	return func(o *options) {
		if port > 0 && port < 65535 {
			o.customPort = port
		}
	}
}

func WithNoTLS() MailOptions {
	return func(o *options) {
		o.noTLS = true
	}
}

func WithHELO(helo string) MailOptions {
	return func(o *options) {
		o.helo = helo
	}
}

func WithCustomUsername(username string) MailOptions {
	return func(o *options) {
		o.username = username
	}
}

func (m *mailSender) Send(sendTo, subject, msg string) ([]string, error) {
	message := mail.NewMsg()
	if err := message.From(m.email); err != nil {
		log.Println("failed to set From address:", m.email)
		return nil, err
	}
	if err := message.To(sendTo); err != nil {
		log.Println("failed to set To address:", sendTo)
		return nil, err
	}
	message.Subject(subject)
	message.SetBodyString(mail.TypeTextPlain, msg)

	client, err := mail.NewClient(m.smtpServer)
	if err != nil {
		log.Println("failed to create mail client")
		return nil, err
	}

	client.SetUsername(m.opts.username)
	client.SetPassword(m.password)

	if m.opts.customAuth {
		client.SetSMTPAuth(mail.SMTPAuthLogin)
		client.SetSSLPort(m.opts.noTLS, true)
		if m.opts.noTLS {
			client.SetTLSConfig(&tls.Config{InsecureSkipVerify: true})
		}
	} else {
		client.SetSMTPAuth(mail.SMTPAuthAutoDiscover)
	}
	if err := client.DialAndSend(message); err != nil {
		log.Println("failed to send mail")
		return nil, err
	}

	messageID := message.GetGenHeader(mail.HeaderMessageID)

	return messageID, nil
}
