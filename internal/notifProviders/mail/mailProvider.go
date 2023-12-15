package mail

import (
	"MailService/internal/db"
	"crypto/tls"
	"github.com/spf13/viper"
	"gopkg.in/gomail.v2"
	"strconv"
)

type Smtp struct {
	SMTPServer   string
	SMTPPort     int
	SMTPUsername string
	SMTPPassword string
	SenderEmail  string
}

func getSmtpConfig() (Smtp, error) {
	smtpConfig := viper.GetStringMapString("smtp")

	smtpPort, err := strconv.Atoi(smtpConfig["port"]) // Viper port değerini string getiriyor ondan inte dönüştürdük
	if err != nil {
		return Smtp{}, err
	}

	return Smtp{
		SMTPServer:   smtpConfig["host"],
		SMTPPort:     smtpPort,
		SMTPUsername: smtpConfig["username"],
		SMTPPassword: smtpConfig["password"],
		SenderEmail:  smtpConfig["sender_email"], //viper camelCase okumuyor :(
	}, nil
}

func SendMail(toEmail string, templateAlias string, siteId int, customData map[string]interface{}) error {
	// Smtp bilgilerini al
	SmtpConfig, err := getSmtpConfig()
	if err != nil {
		return err
	}

	mailRecord, err := db.GetMailContent(templateAlias, siteId, customData)
	if err != nil {
		return err
	}

	// Gönderen bilgileri
	sender := gomail.NewMessage()
	sender.SetHeader("From", SmtpConfig.SenderEmail)
	sender.SetHeader("To", toEmail)
	sender.SetHeader("Subject", mailRecord.Subject)
	sender.SetBody("text/html", mailRecord.Content)

	// SMTP ayarları
	dialer := gomail.NewDialer(SmtpConfig.SMTPServer, SmtpConfig.SMTPPort, SmtpConfig.SMTPUsername, SmtpConfig.SMTPPassword)
	dialer.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	// E-postayı gönder
	if err := dialer.DialAndSend(sender); err != nil {
		return err
	}

	return nil
}
