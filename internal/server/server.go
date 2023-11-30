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
	ToEmail    string                 `json:"email"`
	TemplateID int                    `json:"templateID"`
	IsQueue    bool                   `json:"isQueue"`
	CustomData map[string]interface{} `json:"customData"`
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
			err := mailSender.SendMail(payload.ToEmail, payload.TemplateID, payload.CustomData)
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
	if payload.TemplateID == 0 {
		return errors.New("Template Id Geçersiz")
	}

	_, err := mail.ParseAddress(payload.ToEmail)
	if err != nil {
		return errors.New("Email adresi geçersiz")
	}
	return nil
}
