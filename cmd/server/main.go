package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/onunkwor/flypro-assestment-v2/internal/config"
	"github.com/onunkwor/flypro-assestment-v2/internal/routes"
)

func init() {
	config.LoadEnv()
	err := config.ConnectDatabase()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	config.ConnectRedis()
}

func main() {
	router := gin.Default()
	router.Use(gin.Recovery())
	routes.RegisterUserRoutes(router)
	routes.RegisterExpenseRoutes(router)
	routes.RegisterReportRoutes(router)
	port, err := config.Getenv("PORT")
	if err != nil {
		log.Fatal("Failed to get PORT:", err)
	}
	log.Printf("🚀 Server running on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
