package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/onunkwor/flypro-assestment-v2/internal/handlers"
)

func RegisterExpenseRoutes(router *gin.Engine) {
	expenseHandler := handlers.NewExpenseHandler(nil)
	expenseGroup := router.Group("api/expenses")
	{
		expenseGroup.POST("/", expenseHandler.CreateExpense)
		expenseGroup.GET("/:id")

	}
}
