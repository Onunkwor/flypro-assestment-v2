package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/onunkwor/flypro-assestment-v2/internal/config"
)

func init() {
	config.LoadEnv()
	err := config.ConnectDatabase()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
}

func main() {
	router := gin.Default()
	port, err := config.Getenv("PORT")
	if err != nil {
		log.Fatal("Failed to get PORT:", err)
	}
	log.Printf("🚀 Server running on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
