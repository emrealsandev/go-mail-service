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

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		server.StartServer()
	}()

	// ProcessQueue fonksiyonunu başlat (goroutine olarak)
	wg.Add(1)
	go func() {
		defer wg.Done()
		queue.ProcessQueue("mailQueue")
	}()

	// Ana programın diğer işlemleri burada başlatılır
	// ...

	// Tüm goroutine'ların tamamlanmasını bekleyelim
	wg.Wait()
}
