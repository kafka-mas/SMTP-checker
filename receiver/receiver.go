package receiver

import (
	"fmt"
	"time"

	"github.com/emersion/go-imap/v2"
	"github.com/emersion/go-imap/v2/imapclient"
)

type MailReceiver interface {
	/*
		targetInbox - target mailbox in the mail client, for example "INBOX"

		messageID - the MessageID field in email headers

		Returns true if message found, false otherwise
	*/
	Receive(targetInbox, messageID string) (bool, error)
}

type mailReceiver struct {
	username   string
	password   string
	imapServer string
	opts       options
}

// Constructor; returns MailReceiver
func NewMailReceiver(login, password, imapServer string, opts ...IMAPoptions) MailReceiver {
	mr := &mailReceiver{
		username:   login,
		password:   password,
		imapServer: imapServer,
		opts: options{
			useTLS:         true,
			port:           993,
			customUsername: "",
			timeout:        30,
			pollingTimeout: 5,
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
	pollingTimeout time.Duration
}

// Setup IMAP port
func WithCustomPort(port int) IMAPoptions {
	return func(mr *options) {
		if port > 0 && port < 65535 {
			mr.port = port
		}
	}
}

// Use an insecure connection without verifying the certificate
func WithNoTLS() IMAPoptions {
	return func(o *options) {
		o.useTLS = false
		if o.port == 993 {
			o.port = 143
		}
	}
}

// set username values if it is not specified in the email@domain format
func WithCustomUsername(username string) IMAPoptions {
	return func(o *options) { o.customUsername = username }
}

// Set polling time in seconds
func WithTimeout(timeout time.Duration) IMAPoptions {
	return func(o *options) { o.timeout = timeout }
}

// Set polling timeout in seconds
func WithPollingTimeout(pollingTimeout time.Duration) IMAPoptions {
	return func(o *options) {
		if pollingTimeout > 0 && pollingTimeout <= o.timeout {
			o.pollingTimeout = pollingTimeout
		}
	}
}

func (m *mailReceiver) Receive(targetInbox, messageID string) (bool, error) {
	if messageID == "" {
		return false, fmt.Errorf("messageID is empty string")
	}
	url := fmt.Sprintf("%v:%d", m.imapServer, m.opts.port)
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

	// Auth
	if m.opts.customUsername != "" {
		m.username = m.opts.customUsername
	}
	if err := client.Login(m.username, m.password).Wait(); err != nil {
		return false, fmt.Errorf("failed to login: %w", err)
	}

	_, err = client.Select(targetInbox, nil).Wait()
	if err != nil {
		return false, fmt.Errorf("failed to select inbox: %w", err)
	}

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
		if time.Since(timeStart) >= m.opts.timeout*time.Second {
			return false, nil
		}

		searchRes, err := client.UIDSearch(criteria, nil).Wait()
		if err != nil {
			return false, err
		}
		uids := searchRes.AllUIDs()
		if len(uids) != 0 {
			return true, nil
		}
		time.Sleep(m.opts.pollingTimeout * time.Second)
	}
}
