package email

import (
	"fmt"
	"net/smtp"
	"os"

	"github.com/joho/godotenv"
)

var (
	smtpHost = "smtp.gmail.com"
	smtpPort = "587"
)

func SendEmail(to, subject, body string) error {
	// Загрузить переменные окружения
	err := godotenv.Load()
	if err != nil {
		return fmt.Errorf("Failed to load .env: %w", err)
	}

	senderEmail := os.Getenv("SMTP_EMAIL")
	senderPass := os.Getenv("SMTP_PASS")

	auth := smtp.PlainAuth("", senderEmail, senderPass, smtpHost)
	message := []byte(fmt.Sprintf("Subject: %s\r\n\r\n%s", subject, body))
	addr := smtpHost + ":" + smtpPort

	err = smtp.SendMail(addr, auth, senderEmail, []string{to}, message)
	if err != nil {
		return fmt.Errorf("send mail failed: %w", err)
	}

	return nil
}
