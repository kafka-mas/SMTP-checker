package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"

	"smtp-check/receiver"
	"smtp-check/sender"
)

func main() {
	var envpath string = ".env"

	args := os.Args
	if len(args) > 1 {
		envpath = args[1]
	}


	if err := godotenv.Load(envpath); err != nil {
		log.Fatal("error load .env")
	}

	externalEmail := os.Getenv("EXT_EMAIL")
	externalPass := os.Getenv("EXT_PASS")
	externalImapServer := os.Getenv("EXT_IMAP_SERVER")
	externalSmtpServer := os.Getenv("EXT_SMTP_SERVER")

	internalEmail := os.Getenv("IN_EMAIL")
	internalPass := os.Getenv("IN_PASS")
	internalImapServer := os.Getenv("IN_IMAP_SERVER")
	internalImapPort_s := os.Getenv("IN_IMAP_PORT")
	internalSmtpServer := os.Getenv("IN_SMTP_SERVER")
	internalSmtpPort_s := os.Getenv("IN_SMTP_PORT")
	internalUser := os.Getenv("IN_USER")

	internalImapPort, err := strconv.Atoi(internalImapPort_s)
	if err != nil {
		log.Fatalf("error, param INTERNAL_IMAP_PORT must be integer: %v", err)
	}
	internalSmtpPort, err := strconv.Atoi(internalSmtpPort_s)
	if err != nil {
		log.Fatalf("error, param INTERNAL_IMAP_PORT must be integer: %v", err)
	}

	internalSender := sender.NewMailSender(internalEmail, internalPass, internalSmtpServer,
		sender.WithCustomPort(internalSmtpPort),
		sender.WithCustomUsername(internalUser),
		sender.WithNoTLS(),
	)
	externalSender := sender.NewMailSender(externalEmail, externalPass, externalSmtpServer)

	inToExtID, err := internalSender.Send(externalEmail, "test from internal", "message")
	if err != nil {
		log.Fatalf("error send message: %s", err)
	}
	extToInID, err := externalSender.Send(internalEmail, "test from external", "message")
	if err != nil {
		log.Fatalf("error send message: %s", err)
	}

	internalReceiver := receiver.NewMailReceiver(internalEmail, internalPass, internalImapServer,
		receiver.WithCustomPort(internalImapPort),
		receiver.WithCustomUsername(internalUser),
		receiver.WithNoTLS(),
	)
	externalReceiver := receiver.NewMailReceiver(externalEmail, externalPass, externalImapServer)

	isExtToInFind, err := internalReceiver.Receive("INBOX", extToInID)
	if err != nil {
		log.Fatalln("Error:", err)
	}

	isInToExtFind, err := externalReceiver.Receive("INBOX", inToExtID)
	if err != nil {
		log.Fatalln("Error:", err)
	}

	type InternalToExternal struct {
		MessageID       string `json:"messageID"`
		MessageReceived bool   `json:"messageReceived"`
	}
	type ExternalToInternal struct {
		MessageID       string `json:"messageID"`
		MessageReceived bool   `json:"messageReceived"`
	}
	type Result struct {
		InternalToExternal InternalToExternal `json:"internalToExternal"`
		ExternalToInternal ExternalToInternal `json:"externalToInternal"`
	}

	res := Result{
		InternalToExternal: InternalToExternal{
			MessageID:       inToExtID,
			MessageReceived: isInToExtFind,
		},
		ExternalToInternal: ExternalToInternal{
			MessageID:       extToInID,
			MessageReceived: isExtToInFind,
		},
	}
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false) 
	err = enc.Encode(res)
	if err != nil {
		log.Fatalf("error encoding struct to json: %v", err)
	}

	resJSON := bytes.TrimSpace(buf.Bytes()) 
	fmt.Println(string(resJSON))
}
