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

type EmailDetailResponse struct {
	Success  bool   `json:"success"`
	Email    string `json:"email"`
	IsQueued bool   `json:"is_queued"`
	Message  string `json:"message"`
}

type StatusResponse struct {
	Success bool                  `json:"success"`
	Code    int                   `json:"code"`
	Message string                `json:"message"`
	Emails  []EmailDetailResponse `json:"emails"`
}

func StartServer() {
	router := fiber.New()

	router.Post("/sendMail", func(c *fiber.Ctx) error {
		payload := BulkPayload{}
		var emailDetails []EmailDetailResponse

		// Requesti al
		if err := c.BodyParser(&payload); err != nil {
			return c.Status(400).JSON(createResponse(400, err.Error(), false, emailDetails))
		}

		for _, mailPayload := range payload.Emails {
			var message string
			var success = true

			// Request Validate et
			err := checkBulkMailSendable(&mailPayload)
			if err != nil {
				success = false
				message = err.Error()
			}

			if success != false {
				if !(mailPayload.IsQueue) {
					err := mailSender.SendMail(mailPayload.ToEmail, mailPayload.TemplateAlias, mailPayload.SiteID, mailPayload.CustomData)
					if err != nil {
						success = false
						message = "Mail Gönderimi Başarısız"
					} else {
						success = true
						message = "Mail Gönderimi Başarılı"
						log.Println("Başarılı gönderim:", mailPayload.ToEmail)
					}
				} else {
					err := queue.AddToQueue(mailPayload.ToEmail, mailPayload.TemplateAlias, mailPayload.SiteID, mailPayload.CustomData, "mailQueue")
					if err != nil {
						success = false
						message = "Kuyruğa ekleme işlemi başarısız"
					} else {
						success = true
						message = "Kuyruğa ekleme başarılı"
						log.Println("Kuyruğa ekleme başarılı:", mailPayload.ToEmail)
					}
				}
			}

			emailDetails = append(emailDetails, EmailDetailResponse{
				Success:  success,
				Email:    mailPayload.ToEmail,
				IsQueued: mailPayload.IsQueue,
				Message:  message,
			})
		}
		response := createResponse(200, "işlem tamamlandı", true, emailDetails)
		return c.Status(200).JSON(response)
	})

	err := router.Listen(":3000")
	if err != nil {
		log.Fatal(err)
	}
}

// bunun public olmasına gerek yok
func checkBulkMailSendable(payload *Payload) error {
	if payload.TemplateAlias == "" {
		return errors.New("template Alias geçersiz")
	}

	if payload.SiteID <= 0 && payload.SiteID == 0 {
		return errors.New("siteId geçersiz")
	}

	// E-posta adresini kontrol et
	if _, err := mail.ParseAddress(payload.ToEmail); err != nil {
		return errors.New("email adresi geçersiz")
	}
	return nil
}

// public olmasına gerek yok
func createResponse(code int, message string, success bool, emailDetail []EmailDetailResponse) StatusResponse {
	return StatusResponse{
		Success: success,
		Code:    code,
		Message: message,
		Emails:  emailDetail,
	}
}
