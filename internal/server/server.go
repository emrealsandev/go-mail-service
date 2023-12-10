package server

import (
	mailSender "MailService/internal/notifProviders/mail"
	"MailService/internal/queue"
	"errors"
	"github.com/gofiber/fiber/v2"
	"log"
	"net/mail"
)

type Payload struct {
	// Json struct değerlerinin baş harfi büyük harf ile başlamalıymış
	ToEmail       string                 `json:"email"`
	TemplateAlias string                 `json:"template_alias"`
	SiteID        int                    `json:"site_id"`
	IsQueue       bool                   `json:"is_queue"`
	CustomData    map[string]interface{} `json:"custom_data"`
}

type BulkPayload struct {
	Emails []Payload `json:"emails"`
}

func StartServer() {

	router := fiber.New()

	router.Post("/sendMail", func(c *fiber.Ctx) error {
		payload := BulkPayload{}

		// Requesti al
		if err := c.BodyParser(&payload); err != nil {
			log.Fatal(err)
		}

		// Request Validate et
		err := checkBulkMailSendable(&payload)
		if err != nil {
			log.Fatal(err)
		}

		for _, mailPayload := range payload.Emails {
			if !(mailPayload.IsQueue) {
				err := mailSender.SendMail(mailPayload.ToEmail, mailPayload.TemplateAlias, mailPayload.SiteID, mailPayload.CustomData)
				if err != nil {
					log.Println("Başarısız (kuyruk):", mailPayload.ToEmail)
				} else {
					log.Println("Başarılı (kuyruk):", mailPayload.ToEmail)
				}
			} else {
				err := queue.AddToQueue(mailPayload.ToEmail, mailPayload.TemplateAlias, mailPayload.SiteID, mailPayload.CustomData, "mailQueue")
				if err != nil {
					log.Println("Kuyruğa ekleme başarısız:", mailPayload.ToEmail)
				} else {
					log.Println("Kuyruğa ekleme başarılı:", mailPayload.ToEmail)
				}
			}
		}
		return c.SendString("işlem tamamlandı")
	})

	err := router.Listen(":3000")
	if err != nil {
		log.Fatal(err)
	}
}

// bunun public olmasına gerek yok
func checkBulkMailSendable(payload *BulkPayload) error {
	for _, mailPayload := range payload.Emails {
		if mailPayload.TemplateAlias == "" {
			return errors.New("template Alias geçersiz")
		}

		if mailPayload.SiteID <= 0 {
			return errors.New("siteId geçersiz")
		}

		// E-posta adresini kontrol et
		if _, err := mail.ParseAddress(mailPayload.ToEmail); err != nil {
			return errors.New("email adresi geçersiz")
		}
	}
	return nil
}
