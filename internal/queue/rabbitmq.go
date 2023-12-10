// internal/queue/rabbitmq.go

package queue

import (
	"fmt"
	"github.com/spf13/viper"
	"github.com/streadway/amqp"
	"log"
)

var (
	conn *amqp.Connection
	ch   *amqp.Channel
)

// InitRabbitMQConnection RabbitMQ bağlantısını başlatır
func InitRabbitMQConnection() {
	var err error
	rabbitMQConfig := viper.GetStringMapString("rabbitmq")

	// RabbitMQ bağlantısı yap
	conn, err = amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%s/",
		rabbitMQConfig["username"],
		rabbitMQConfig["password"],
		rabbitMQConfig["host"],
		rabbitMQConfig["port"],
	))
	if err != nil {
		log.Fatalf("RabbitMQ'ya bağlanırken hata oluştu: %v", err)
	}

	// Kanal oluştur
	ch, err = conn.Channel()
	if err != nil {
		log.Fatalf("RabbitMQ kanalı oluşturulurken hata oluştu: %v", err)
	}

}

// CloseRabbitMQConnection RabbitMQ bağlantısını kapatır
func CloseRabbitMQConnection() {
	if ch != nil {
		err := ch.Close()
		if err != nil {
			log.Fatalf("RabbitMQ bağlantısı kapatılamadı: %v", err)
		}
	}
	if conn != nil {
		err := conn.Close()
		if err != nil {
			log.Fatalf("RabbitMQ bağlantısı kapatılamadı: %v", err)
		}
	}
}

// PublishToQueue RabbitMQ kuyruğuna mesaj ekler
func PublishToQueue(queueName string, message string) error {
	err := ch.Publish(
		"",        // exchange
		queueName, // routing key (kuyruk adı)
		false,     // mandatory
		false,     // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        []byte(message),
		},
	)
	if err != nil {
		return fmt.Errorf("Kuyruğa mesaj eklenirken hata oluştu: %v", err)
	}

	return nil
}

func ConsumeQueue(queueName string) (<-chan amqp.Delivery, error) {
	messages, err := ch.Consume(
		queueName, // queue
		"",        // consumer
		true,      // auto-ack
		false,     // exclusive
		false,     // no-local
		false,     // no-wait
		nil,       // args
	)
	if err != nil {
		return nil, err
	}

	return messages, nil
}
