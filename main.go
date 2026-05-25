package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"

	"smtp-check/receiver"
	_ "smtp-check/sender"
)

func main() {
	fmt.Println("Hello")
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("Ошибка загрузки файла .env")
	}

	// externalMail := os.Getenv("SEND_EXT")
	// externalPass := os.Getenv("PASS_EXT")
	// internalMail := os.Getenv("SEND_IN")
	internalPass := os.Getenv("PASS_IN")
	// bolidMailer := sender.NewMailSender("checkmail", internalPass, "mail-serv.bolid.ru")
	// bolidMailer := sender.NewMailSender(internalMail, internalPass, "mail-serv.bolid.ru",
	// 	sender.WithCustomAuth(),
	// 	sender.WithCustomUsername("checkmail"),
	// 	sender.WithCustomPort(25),
	// 	sender.WithNoTLS(),
	// 	// sender.WithHELO("mail-serv.bolid.ru"),
	// )
	//
	// mailMailer := sender.NewMailSender(externalMail, externalPass, "smtp.mail.ru")
	//
	// var msgBolidID []string
	// msgBolidID, err := bolidMailer.Send(externalMail, "Test mail", "Message")
	// if err != nil {
	// 	log.Fatalf("error send message: %s", err)
	// }
	//
	// fmt.Println("Sent message ID:", msgBolidID[0])
	//
	// var msgMailID []string
	// msgMailID, err = mailMailer.Send(internalMail, "Test mail", "Message")
	// if err != nil {
	// 	log.Fatalf("error send message: %s", err)
	// }
	//
	// fmt.Println("Sent message ID:", msgMailID[0])


	// var msg []string = []string{"123"}
	// err = receiver.WaitForEmail(externalMail, externalPass, "INBOX", msg[0])
	imap := receiver.NewMailReceiver("", internalPass, "")
	err := imap.Receive("INBOX","")
	if err != nil {
		log.Fatalln("Error:", err)
	}
}
