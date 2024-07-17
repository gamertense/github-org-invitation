package main

import (
	"log"
	"my-app/internal/handler"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	router := gin.Default()

	router.POST("/send-invitation", handler.SendInvitation)

	router.Run(":8080") // Listen and serve on 0.0.0.0:8080
}
