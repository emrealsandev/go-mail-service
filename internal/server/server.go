package server

import (
	mailSender "MailService/internal/notifProviders/mail"
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

func StartServer() {

	router := fiber.New()

	router.Post("/sendMail", func(c *fiber.Ctx) error {
		payload := Payload{}

		// Requesti al
		// TODO: Bu bulk alıp işleyecek şekilde düzenlenebilir
		if err := c.BodyParser(&payload); err != nil {
			log.Fatal(err)
		}

		// Request Validate et
		err := checkMailSendable(&payload)
		if err != nil {
			log.Fatal(err)
		}

		// TODO: Kuyruklu gönderim ?
		if payload.IsQueue == false {
			err := mailSender.SendMail(payload.ToEmail, payload.TemplateAlias, payload.SiteID, payload.CustomData)
			if err != nil {
				return c.SendString("Başarısız")
			} else {
				return c.SendString("Başarılı")
			}
		} else {
			return c.SendString("Kuyruklu Gönderim Şuan aktif değildir")
		}
	})

	err := router.Listen(":3000")
	if err != nil {
		log.Fatal(err)
	}
}

// bunun public olmasına gerek yok
func checkMailSendable(payload *Payload) error {
	if payload.TemplateAlias == "" {
		return errors.New("template Alias geçersiz")
	}

	if payload.SiteID <= 0 {
		return errors.New("siteId geçersiz")
	}

	_, err := mail.ParseAddress(payload.ToEmail)
	if err != nil {
		return errors.New("email adresi geçersiz")
	}
	return nil
}
