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

type StatusResponse struct {
	Success bool   `json:"success"`
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func StartServer() {

	router := fiber.New()

	router.Post("/sendMail", func(c *fiber.Ctx) error {
		successMessage := "İşlem Başarılı"
		payload := BulkPayload{}

		// Requesti al
		if err := c.BodyParser(&payload); err != nil {
			return c.Status(400).JSON(createResponse(400, err.Error(), false))
		}

		// Request Validate et
		err := checkBulkMailSendable(&payload)
		if err != nil {
			return c.Status(400).JSON(createResponse(400, err.Error(), false))
		}

		for _, mailPayload := range payload.Emails {
			if !(mailPayload.IsQueue) {
				err := mailSender.SendMail(mailPayload.ToEmail, mailPayload.TemplateAlias, mailPayload.SiteID, mailPayload.CustomData)
				if err != nil {
					return c.Status(400).JSON(createResponse(400, "Mail Gönderimi Başarısız", false))
				} else {
					log.Println("Başarılı gönderim:", mailPayload.ToEmail)
					successMessage = "Mail gönderimi başarılı"
				}
			} else {
				err := queue.AddToQueue(mailPayload.ToEmail, mailPayload.TemplateAlias, mailPayload.SiteID, mailPayload.CustomData, "mailQueue")
				if err != nil {
					return c.Status(400).JSON(createResponse(400, "Kuyruğa ekleme işlemi başarısız", false))
				} else {
					successMessage = "Kuyruğa ekleme başarılı"
					log.Println("Kuyruğa ekleme başarılı:", mailPayload.ToEmail)
				}
			}
		}
		successResponse := createResponse(200, successMessage, true)
		return c.Status(200).JSON(successResponse)
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

		if mailPayload.SiteID <= 0 && mailPayload.SiteID == 0 {
			return errors.New("siteId geçersiz")
		}

		// E-posta adresini kontrol et
		if _, err := mail.ParseAddress(mailPayload.ToEmail); err != nil {
			return errors.New("email adresi geçersiz")
		}
	}
	return nil
}

// public olmasına gerek yok
func createResponse(code int, message string, success bool) StatusResponse {
	return StatusResponse{
		Success: success,
		Code:    code,
		Message: message,
	}
}
