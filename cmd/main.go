package main

import (
	"MailService/internal/db"
	"MailService/internal/queue"
	"MailService/internal/server"
	"github.com/spf13/viper"
	"log"
	"sync"
)

func main() {
	// Konfigürasyon dosyasını yükle
	viper.AddConfigPath("/Users/emre.alsan/MyPersonalProjects/MailService/configs")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Yapılandırma dosyası okunurken hata oluştu: %v", err)
	}

	// RabitMQ bağlantısını başlat
	queue.InitRabbitMQConnection()
	defer queue.CloseRabbitMQConnection()

	// Veritabanı bağlantısını başlat
	db.InitDB()
	defer db.CloseDB()

	var wg sync.WaitGroup

	// Server başlat
	wg.Add(1)
	go func() {
		defer wg.Done()
		server.StartServer()
	}()

	// ProcessQueue fonksiyonunu başlat
	wg.Add(1)
	go func() {
		defer wg.Done()
		queue.ProcessQueue("mailQueue")
	}()
	wg.Wait()
}
