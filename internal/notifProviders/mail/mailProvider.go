package mail

import (
	"crypto/tls"
	"github.com/spf13/viper"
	"gopkg.in/gomail.v2"
	"strconv"
)

func SendMail(to, subject, body string) error {
	smtpConfig := viper.GetStringMapString("smtp")

	username := smtpConfig["username"]
	password := smtpConfig["password"]
	smtpHost := smtpConfig["host"]
	smtpPort, err := strconv.Atoi(smtpConfig["port"])
	if err != nil {
		return err
	}

	from := "admin@admin.com"

	// Gönderen bilgileri
	sender := gomail.NewMessage()
	sender.SetHeader("From", from)
	sender.SetHeader("To", to)
	sender.SetHeader("Subject", subject)
	sender.SetBody("text/html", body)

	// SMTP ayarları
	dialer := gomail.NewDialer(smtpHost, smtpPort, username, password)
	dialer.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	// E-postayı gönder
	if err := dialer.DialAndSend(sender); err != nil {
		return err
	}

	return nil
}
