package main

import (
	"MailService/internal/db"
	"MailService/internal/queue"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"log"
)

const DefaultQueue = "mailQueue"

func main() {
	viper.AddConfigPath("/Users/emre.alsan/MyPersonalProjects/MailService/configs")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Yapılandırma dosyası okunurken hata oluştu: %v", err)
	}

	// RabitMQ bağlantısını başlat
	queue.InitRabbitMQConnection()
	defer queue.CloseRabbitMQConnection()

	db.InitDB()
	defer db.CloseDB()

	var queueName string
	pflag.StringVarP(&queueName, "queue", "q", "", "Kuyruk adı")
	pflag.Parse()

	if queueName == "" {
		queueName = DefaultQueue
	}

	// RabbitMQ kuyruğunu işle
	queue.ProcessQueue(queueName)
}
