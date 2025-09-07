package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/onunkwor/flypro-assestment-v2/internal/dto"
	"github.com/onunkwor/flypro-assestment-v2/internal/models"
	"github.com/onunkwor/flypro-assestment-v2/internal/repository"
	"github.com/onunkwor/flypro-assestment-v2/internal/services"
	"github.com/onunkwor/flypro-assestment-v2/internal/utils"
)

type ReportHandler interface {
	CreateReport(c *gin.Context)
	AddExpenseToReport(c *gin.Context)
	SubmitReport(c *gin.Context)
	GetReportExpenses(c *gin.Context)
}
type reportHandler struct {
	reportService services.ReportService
}

func NewReportHandler(reportService services.ReportService) ReportHandler {
	return &reportHandler{
		reportService: reportService,
	}
}

func (h *reportHandler) CreateReport(c *gin.Context) {
	var request dto.CreateReportRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		formatted := utils.FormatValidationError(err)
		utils.ValidationErrorResponse(c, formatted)
		return
	}
	request.Sanitize()
	report := models.ExpenseReport{
		UserID: request.UserID,
		Title:  request.Title,
	}
	if err := h.reportService.CreateReport(c.Request.Context(), &report); err != nil {
		utils.InternalServerErrorResponse(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Report created successfully"})
}

func (h *reportHandler) AddExpenseToReport(c *gin.Context) {
	var request dto.AddExpenseToReportRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		formatted := utils.FormatValidationError(err)
		utils.ValidationErrorResponse(c, formatted)
		return
	}

	reportID := c.GetUint("reportID")

	if err := h.reportService.AddExpenseToReport(c.Request.Context(), reportID, request.ExpenseID); err != nil {
		switch err {
		case repository.ErrReportNotFound:
			utils.NotFoundResponse(c, "report not found")
			return
		case repository.ErrExpenseNotFound:
			utils.NotFoundResponse(c, "expense not found")
			return
		default:
			utils.InternalServerErrorResponse(c, err)
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Expense added to report successfully",
	})

}

func (h *reportHandler) SubmitReport(c *gin.Context) {
	reportID := c.GetUint("reportID")

	err := h.reportService.SubmitReport(c.Request.Context(), reportID)
	if err != nil {
		switch err {
		case services.ErrInvalidReportState:
			utils.BadRequestResponse(c, "report cannot be submitted in current state")
			return
		default:
			utils.InternalServerErrorResponse(c, err)
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Report submitted successfully",
	})
}

func (h *reportHandler) GetReportExpenses(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Query("user_id"), 10, 64)
	if err != nil || userID == 0 {
		utils.BadRequestResponse(c, "invalid user ID")
		return
	}

	offset, err := strconv.Atoi(c.DefaultQuery("offset", "0"))
	if err != nil || offset < 0 {
		utils.BadRequestResponse(c, "invalid offset")
		return
	}

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil || limit <= 0 {
		utils.BadRequestResponse(c, "invalid limit")
		return
	}

	reports, err := h.reportService.GetReportExpenses(c.Request.Context(), uint(userID), offset, limit)
	if err != nil {
		utils.InternalServerErrorResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":   reports,
		"count":  len(reports),
		"offset": offset,
		"limit":  limit,
	})
}
