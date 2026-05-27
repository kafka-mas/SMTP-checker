package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"

	"smtp-check/receiver"
	"smtp-check/sender"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
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

	externalMsgID, err := internalSender.Send(externalEmail, "test from internal", "message")
	if err != nil {
		log.Fatalf("error send message: %s", err)
	}
	internalMsgID, err := externalSender.Send(internalEmail, "test from external", "message")
	if err != nil {
		log.Fatalf("error send message: %s", err)
	}

	fmt.Println("Sent internal message ID:", internalMsgID)
	fmt.Println("Send external message ID:", externalMsgID)

	internalReceiver := receiver.NewMailReceiver(internalEmail, internalPass, internalImapServer,
		receiver.WithCustomPort(internalImapPort),
		receiver.WithCustomUsername(internalUser),
		receiver.WithNoTLS(),
	)
	externalReceiver := receiver.NewMailReceiver(externalEmail, externalPass, externalImapServer)

	isInternalMailFind, err := internalReceiver.Receive("INBOX", internalMsgID)
	if err != nil {
		log.Fatalln("Error:", err)
	}

	isExternalMailFind, err := externalReceiver.Receive("INBOX", externalMsgID)
	if err != nil {
		log.Fatalln("Error:", err)
	}

	fmt.Println("Internal mail find:", isInternalMailFind)
	fmt.Println("External mail find:", isExternalMailFind)
}
