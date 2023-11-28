package server

import (
	"MailService/internal/db"
	"MailService/internal/notifProviders/mail"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

func StartServer() {
	router := gin.Default()

	router.GET("/send-mail/:templateID/:to", func(c *gin.Context) {
		templateIDStr := c.Param("templateID")
		to := c.Param("to")

		templateID, err := strconv.Atoi(templateIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid template ID"})
			return
		}

		// Veritabanından mail şablonunu al
		templateContent, err := db.GetMailTemplate(templateID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			return
		}

		// Mail gönderme işlemini gerçekleştir
		err = mail.SendMail(to, "Mail Konusu", templateContent)
		if err != nil {
			fmt.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "success"})
	})

	err := router.Run(":3000")
	if err != nil {
		log.Fatal(err)
	}
}
