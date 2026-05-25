package receiver

import (
	"fmt"
	// "github.com/emersion/go-imap/v2"
	"github.com/emersion/go-imap/v2/imapclient"
)

type MailReceiver interface {
	Receive(targetInbox, messageID string) error
}

type mailReceiver struct {
	username   string
	password   string
	imapServer string
	opts       options
}

func NewMailReceiver(login, password, imapServer string, opts ...IMAPoptions) MailReceiver {
	mr := &mailReceiver{
		username:   login,
		password:   password,
		imapServer: imapServer,
		opts: options{
			useTLS: false,
			port:   143,
		},
	}
	return mr
}

type IMAPoptions func(*mailReceiver)

type options struct {
	port           int
	useTLS         bool
	customUsername string
}

func WithCustomPort(port int) IMAPoptions {
	return func(mr *mailReceiver) {
		if port >0 && port < 65535 {
			mr.opts.port = port
		}
	}
}

func (m *mailReceiver) Receive(targetInbox, messageID string) error {
	// Настройки IMAP для Mail.ru
	fmt.Println("Подключение к серверу...")
	url := fmt.Sprintf("%v:%d", m.imapServer, m.opts.port)
	fmt.Println(url)
	client, err := imapclient.DialInsecure(url, nil)
	if err != nil {
		return fmt.Errorf("failed to dial IMAP: %w", err)
	}
	defer client.Logout()
	fmt.Println("Подключено к серверу")

	// Аутентификация
	if err := client.Login(m.username, m.password).Wait(); err != nil {
		return fmt.Errorf("failed to login: %w", err)
	}
	fmt.Println("Аутентификация успешна")

	mbox, err := client.Select(targetInbox, nil).Wait()
	if err != nil {
		return fmt.Errorf("failed to select inbox: %w", err)
	}

	fmt.Printf("Всего писем: %d\n", mbox.NumMessages)

	if mbox.NumMessages == 0 {
		return nil
	}

	return nil
}
