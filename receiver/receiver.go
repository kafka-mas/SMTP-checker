package receiver

import (
	"fmt"
	"time"

	"github.com/emersion/go-imap/v2"
	"github.com/emersion/go-imap/v2/imapclient"
)

type MailReceiver interface {
	Receive(targetInbox, messageID string) (bool, error)
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
			useTLS:         true,
			port:           993,
			customUsername: "",
			timeout:        time.Duration(30 * time.Second),
		},
	}

	for _, opt := range opts {
		opt(&mr.opts)
	}

	return mr
}

type IMAPoptions func(*options)

type options struct {
	port           int
	useTLS         bool
	customUsername string
	timeout        time.Duration
}

func WithCustomPort(port int) IMAPoptions {
	return func(mr *options) {
		if port > 0 && port < 65535 {
			mr.port = port
		}
	}
}

func WithNoTLS() IMAPoptions {
	return func(o *options) { o.useTLS = false }
}

func WithCustomUsername(username string) IMAPoptions {
	return func(o *options) { o.customUsername = username }
}

//Set polling time in seconds
func SetTimeout(timeout time.Duration) IMAPoptions {
	return func(o *options) { o.timeout = timeout * time.Second }
}

func (m *mailReceiver) Receive(targetInbox, messageID string) (bool, error) {
	fmt.Println("Подключение к серверу...")
	url := fmt.Sprintf("%v:%d", m.imapServer, m.opts.port)
	fmt.Println(url)
	client, err := func(url string) (*imapclient.Client, error) {
		if m.opts.useTLS {
			client, err := imapclient.DialTLS(url, nil)
			return client, err
		} else {
			client, err := imapclient.DialInsecure(url, nil)
			return client, err
		}
	}(url)
	if err != nil {
		return false, fmt.Errorf("failed to dial IMAP: %w", err)
	}
	defer client.Logout()
	fmt.Println("Подключено к серверу")

	// Аутентификация
	if m.opts.customUsername != "" {
		m.username = m.opts.customUsername
	}
	if err := client.Login(m.username, m.password).Wait(); err != nil {
		return false, fmt.Errorf("failed to login: %w", err)
	}
	fmt.Println("Аутентификация успешна")

	mbox, err := client.Select(targetInbox, nil).Wait()
	if err != nil {
		return false, fmt.Errorf("failed to select inbox: %w", err)
	}

	fmt.Printf("Всего писем: %d\n", mbox.NumMessages)

	if mbox.NumMessages == 0 {
		return false, nil
	}

	// 2. Настройка критериев поиска по Header
	criteria := &imap.SearchCriteria{
		Header: []imap.SearchCriteriaHeaderField{
			{
				Key:   "Message-Id",
				Value: messageID,
			},
		},
	}

	timeStart := time.Now()
	for {
		if time.Since(timeStart) >= m.opts.timeout {
			return false, fmt.Errorf("error get message with ID %v", messageID)
		}

		searchRes, err := client.UIDSearch(criteria, nil).Wait()
		if err != nil {
			return false, err
		}
		uids := searchRes.AllUIDs()
		if len(uids) != 0 {
			return true, nil
		}
		time.Sleep(time.Duration(5) * time.Second)
	}
}
