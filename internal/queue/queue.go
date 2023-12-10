// internal/queue/queue.go

package queue

import (
	"MailService/internal/notifProviders/mail"
	"encoding/json"
	"fmt"
	"log"
)

type RequestStruct struct {
	Email         string      `json:"email"`
	TemplateAlias string      `json:"template_alias"`
	SiteID        int         `json:"site_id"`
	CustomData    interface{} `json:"custom_data"`
}

func AddToQueue(toEmail string, templateAlias string, siteId int, customData map[string]interface{}, queueName string) error {
	record, err := prepareQueueRecord(toEmail, templateAlias, siteId, customData)
	if err != nil {
		return err
	}
	err = PublishToQueue(queueName, record)
	if err != nil {
		return err
	}

	return nil
}

// public olmasına gerek yok
func prepareQueueRecord(toEmail string, templateAlias string, siteID int, customData map[string]interface{}) (string, error) {
	request := RequestStruct{
		Email:         toEmail,
		TemplateAlias: templateAlias,
		SiteID:        siteID,
		CustomData:    customData,
	}
	jsonData, err := json.Marshal(request)
	if err != nil {
		// Errorf fonksiyonu error tipi dönüyor
		return "", fmt.Errorf("JSON formatına çevirme hatası: %v", err)
	}
	return string(jsonData), nil
}

func ProcessQueue(queueName string) {
	messages, err := ConsumeRabbitQueue(queueName)
	if err != nil {
		log.Fatalf("Kuyruktan mesajları alırken hata oluştu: %v", err)
	}

	for message := range messages {
		// Kuyruktan gelen mesajı işle
		handleQueueMessage(message.Body)
	}
}

func handleQueueMessage(body []byte) {
	// Mesaj içeriğini JSON formatından çıkar
	var mailRequest RequestStruct
	if err := json.Unmarshal(body, &mailRequest); err != nil {
		log.Fatalf("Kuyruktan gelen mesajı işlerken hata oluştu: %v", err)
		return
	}

	// Kuyruktan gelen her bir mesajı işle
	log.Printf("İşlenen kayıt bilgileri: Email: %s, Template:%s, SiteId:%v, CustomData:%v", mailRequest.Email, mailRequest.TemplateAlias, mailRequest.SiteID, mailRequest.CustomData)
	err := mail.SendMail(mailRequest.Email, mailRequest.TemplateAlias, mailRequest.SiteID, mailRequest.CustomData.(map[string]interface{}))
	if err != nil {
		log.Printf("E-posta gönderirken hata oluştu: %v", err)
	}
}
