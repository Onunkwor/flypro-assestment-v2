package routes

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/onunkwor/flypro-assestment-v2/internal/config"
	"github.com/onunkwor/flypro-assestment-v2/internal/handlers"
	"github.com/onunkwor/flypro-assestment-v2/internal/repository"
	"github.com/onunkwor/flypro-assestment-v2/internal/services"
)

func RegisterExpenseRoutes(router *gin.Engine) {
	currencyApi, err := config.Getenv("CURRENCY_API")
	if err != nil {
		log.Fatal("CURRENCY_API not set in environment")
	}
	expenseRepository := repository.NewExpenseRepository(config.DB)
	currencyService := services.NewCurrencyService(config.Redis, currencyApi, 60*time.Hour)
	expenseService := services.NewExpenseService(config.Redis, currencyService, expenseRepository)
	expenseHandler := handlers.NewExpenseHandler(expenseService)
	expenseGroup := router.Group("api/expenses")
	{
		expenseGroup.POST("/", expenseHandler.CreateExpense)
		expenseGroup.GET("/:id")
		expenseGroup.GET("/", expenseHandler.GetExpenses)
		expenseGroup.PUT("/:id", expenseHandler.UpdateExpense)
		expenseGroup.DELETE("/:id", expenseHandler.DeleteExpense)
	}
}
