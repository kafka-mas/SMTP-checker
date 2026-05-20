package sender

import (
	"log"

	"github.com/wneessen/go-mail"
)

type mailSender struct {
	SendFrom string
	Password string
}

type Mail interface {
	Send(sendTo, subject, msg string) error
}

func NewMailSender(from, password string) Mail {
	return &mailSender{
		SendFrom: from,
		Password: password,
	}
}

func (m mailSender) Send(sendTo, subject, msg string) error {
	message := mail.NewMsg()
	if err := message.From(m.SendFrom); err != nil {
		log.Fatalf("failed to set From address: %s", err)
		return err
	}
	if err := message.To(sendTo); err != nil {
		log.Fatalf("failed to set To address: %s", err)
		return err
	}
	message.Subject(subject)
	message.SetBodyString(mail.TypeTextPlain, msg)
	client, err := mail.NewClient("smtp.gmail.com", mail.WithSMTPAuth(mail.SMTPAuthAutoDiscover),
		mail.WithUsername(m.SendFrom), mail.WithPassword(m.Password))
	if err != nil {
		log.Fatalf("failed to create mail client: %s", err)
	}
	if err := client.DialAndSend(message); err != nil {
		log.Fatalf("failed to send mail: %s", err)
	}
	return nil
}
