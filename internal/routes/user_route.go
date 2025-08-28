package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/onunkwor/flypro-assestment-v2/internal/config"
	"github.com/onunkwor/flypro-assestment-v2/internal/handlers"
	"github.com/onunkwor/flypro-assestment-v2/internal/repository"
	"github.com/onunkwor/flypro-assestment-v2/internal/services"
)

func RegisterUserRoutes(router *gin.Engine) {
	userRepo := repository.NewUserRepository(config.DB)
	userService := services.NewUserService(config.Redis, userRepo)
	userHandler := handlers.NewUserHandler(userService)
	userGroup := router.Group("/api/users")
	{
		userGroup.POST("/", userHandler.CreateUser)
		userGroup.GET("/:id", userHandler.GetUserByID)

	}
}
