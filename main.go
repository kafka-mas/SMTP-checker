package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"

	"smtp-check/sender"
)

func main() {
	fmt.Println("Hello")
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("Ошибка загрузки файла .env")
	}

	sendFrom := os.Getenv("SEND_FROM")
	sendPass := os.Getenv("PASS")
	sendTo := os.Getenv("SEND_TO")

	m := sender.NewMailSender(sendFrom, sendPass)

	if err := m.Send(sendTo, "Test mail", "Message"); err != nil {
		log.Fatalf("error send message: %s", err)
	}
}
