package receiver

import (
	"fmt"
	// "github.com/emersion/go-imap/v2"
	"github.com/emersion/go-imap/v2/imapclient"
)

func WaitForEmail(email, pass, targetInbox, messageID string) error {
	// Настройки IMAP для Mail.ru
	fmt.Println("Подключение к серверу...")
	client, err := imapclient.DialTLS("imap.mail.ru:993", nil)
	if err != nil {
		return fmt.Errorf("failed to dial IMAP: %w", err)
	}
	defer client.Logout()
	fmt.Println("Подключено к серверу")

	// Аутентификация
	if err := client.Login(email, pass).Wait(); err != nil {
		return fmt.Errorf("failed to login: %w", err)
	}
	fmt.Println("Аутентификация успешна")

	// Выбираем нужный почтовый ящик (INBOX)
	// mbox, err := client.Select(targetInbox, nil).Wait()
	// if err != nil {
	// 	return fmt.Errorf("failed to select inbox: %w", err)
	// }
	// fmt.Printf("Всего писем: %d\n", mbox.NumMessages)
	//
	// // Если писем нет, завершаем работу
	// if mbox.Messages == 0 {
	// 	return
	// }

	mbox, err := client.Select(targetInbox, nil).Wait()
	if err != nil{
		return fmt.Errorf("failed to select inbox: %w", err)
	}
	
	fmt.Printf("Всего писем: %d\n", mbox.NumMessages)

	if mbox.NumMessages == 0 {
		return nil
	}

	return nil
}
