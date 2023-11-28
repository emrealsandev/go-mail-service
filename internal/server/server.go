package server

import (
	"MailService/internal/notifProviders/mail"
	"github.com/gofiber/fiber/v2"
	"log"
)

type Payload struct {
	// Json struct değerlerinin baş harfi büyük harf ile başlamalıymış
	Name       string `json:"name"`
	Email      string `json:"email"`
	TemplateID int    `json:"templateID"`
	OrderID    string `json:"orderID"`
	IsQueue    bool   `json:"isQueue"`
}

func StartServer() {

	router := fiber.New()

	router.Post("/sendMail", func(c *fiber.Ctx) error {
		payload := Payload{}

		if err := c.BodyParser(&payload); err != nil {
			log.Fatal(err)
		}

		if payload.IsQueue == false {
			err := mail.SendMail(payload.Email, "Başlık", "Mail contenti böyle olacak gibi bir şey")
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
