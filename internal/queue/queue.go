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

// EnqueueRequestToQueue gelen isteği kuyruğa ekler
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

func prepareQueueRecord(toEmail string, templateAlias string, siteID int, customData map[string]interface{}) (string, error) {
	request := RequestStruct{
		Email:         toEmail,
		TemplateAlias: templateAlias,
		SiteID:        siteID,
		CustomData:    customData,
	}
	jsonData, err := json.Marshal(request)
	if err != nil {
		return "", fmt.Errorf("JSON formatına çevirme hatası: %v", err)
	}
	return string(jsonData), nil
}

func ProcessQueue(queueName string) {
	messages, err := ConsumeQueue(queueName)
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
		log.Printf("Kuyruktan gelen mesajı işlerken hata oluştu: %v", err)
		return
	}

	// JSON formatındaki customData'yı map[string]interface{} türüne çözümle
	var customData map[string]interface{}
	switch v := mailRequest.CustomData.(type) {
	case string:
		// Zaten bir JSON stringi ise doğrudan çözümle
		err := json.Unmarshal([]byte(v), &customData)
		if err != nil {
			log.Printf("CustomData çözümlenirken hata oluştu: %v", err)
			return
		}
	case map[string]interface{}:
		// Zaten bir harita ise doğrudan kullan
		customData = v
	default:
		log.Printf("CustomData JSON formatında değil: %v", mailRequest.CustomData)
		return
	}

	// Kuyruktan gelen her bir mesajı işle
	// Mesajı mailProvider.go'daki SendMail fonksiyonuna yönlendir
	err := mail.SendMail(mailRequest.Email, mailRequest.TemplateAlias, mailRequest.SiteID, customData)
	if err != nil {
		log.Printf("E-posta gönderirken hata oluştu: %v", err)
	}
}
