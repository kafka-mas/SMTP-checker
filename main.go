package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"

	"smtp-check/sender"
	"smtp-check/receiver"
)

func main() {
	fmt.Println("Hello")
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("Ошибка загрузки файла .env")
	}

	externalMail := os.Getenv("SEND_EXT")
	externalPass := os.Getenv("PASS_EXT")
	internalMail := os.Getenv("SEND_IN")
	internalPass := os.Getenv("PASS_IN")
	m := sender.NewMailSender(externalMail, externalPass, "smtp.mail.ru")

	var msgID []string
	msgID_p, err := m.Send(internalMail, "Test mail", "Message")
	if err != nil {
		log.Fatalf("error send message: %s", err)
	}
	msgID = *msgID_p

	fmt.Println("Sent message ID:", msgID[0])
	
	err = receiver.WaitForEmail(internalMail, internalPass, "INBOX", msgID[0])
}
