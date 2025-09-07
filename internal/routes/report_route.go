package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/onunkwor/flypro-assestment-v2/internal/config"
	"github.com/onunkwor/flypro-assestment-v2/internal/handlers"
	"github.com/onunkwor/flypro-assestment-v2/internal/middleware"
	"github.com/onunkwor/flypro-assestment-v2/internal/repository"
	"github.com/onunkwor/flypro-assestment-v2/internal/services"
)

func RegisterReportRoutes(router *gin.Engine) {
	reportRepository := repository.NewReportRepository(config.DB)
	expenseRepository := repository.NewExpenseRepository(config.DB)
	userRepository := repository.NewUserRepository(config.DB)

	reportService := services.NewReportService(
		reportRepository,
		expenseRepository,
		userRepository,
		config.Redis,
	)

	reportHandler := handlers.NewReportHandler(reportService)

	reportRoutes := router.Group("/api/reports")
	{
		reportRoutes.POST("/", reportHandler.CreateReport)
		reportRoutes.POST(
			"/:id/expenses",
			middleware.ReportOwnershipMiddleware(reportRepository),
			middleware.ExpenseOwnershipMiddleware(expenseRepository),
			reportHandler.AddExpenseToReport,
		)
		reportRoutes.PUT(
			"/:id/submit",
			middleware.ReportOwnershipMiddleware(reportRepository),
			reportHandler.SubmitReport,
		)
		reportRoutes.GET("/", reportHandler.GetReportExpenses)
	}
}
