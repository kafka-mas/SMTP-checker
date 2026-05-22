package sender

import (
	"log"

	"github.com/wneessen/go-mail"
)

type mailSender struct {
	SendFrom   string
	Password   string
	smtpServer string
}

type Mail interface {
	Send(sendTo, subject, msg string) (*[]string, error)
}

func NewMailSender(from, password, smtpServer string) Mail {
	return &mailSender{
		SendFrom:   from,
		Password:   password,
		smtpServer: smtpServer,
	}
}

func (m mailSender) Send(sendTo, subject, msg string) (*[]string, error) {
	message := mail.NewMsg()
	if err := message.From(m.SendFrom); err != nil {
		log.Println("failed to set From address:", m.SendFrom)
		return nil, err
	}
	if err := message.To(sendTo); err != nil {
		log.Println("failed to set To address:", sendTo)
		return nil, err
	}
	message.Subject(subject)
	message.SetBodyString(mail.TypeTextPlain, msg)
	client, err := mail.NewClient(m.smtpServer, mail.WithSMTPAuth(mail.SMTPAuthAutoDiscover),
		mail.WithUsername(m.SendFrom), mail.WithPassword(m.Password))
	if err != nil {
		log.Println("failed to create mail client")
		return nil, err
	}
	if err := client.DialAndSend(message); err != nil {
		log.Println("failed to send mail")
		return nil, err
	}

	messageID := message.GetGenHeader(mail.HeaderMessageID)

	return &messageID, nil
}
