package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/onunkwor/flypro-assestment-v2/internal/dto"
	"github.com/onunkwor/flypro-assestment-v2/internal/models"
	"github.com/onunkwor/flypro-assestment-v2/internal/services"
	"github.com/onunkwor/flypro-assestment-v2/internal/utils"
	"gorm.io/gorm"
)

type ExpenseHandler struct {
	service services.ExpenseService
}

func NewExpenseHandler(service services.ExpenseService) *ExpenseHandler {
	return &ExpenseHandler{service: service}
}

func (h *ExpenseHandler) CreateExpense(c *gin.Context) {
	var request dto.CreateExpenseRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		formatted := utils.FormatValidationError(err)
		utils.ValidationErrorResponse(c, formatted)
		return
	}
	exp := &models.Expense{
		UserID:   request.UserId,
		Amount:   request.Amount,
		Currency: request.Currency,
	}
	if err := h.service.CreateExpense(c.Request.Context(), exp); err != nil {
		utils.InternalServerErrorResponse(c, err)
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Expense created successfully"})
}

func (h *ExpenseHandler) GetExpenseByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		utils.BadRequestResponse(c, "Invalid expense ID")
		return
	}

	expense, err := h.service.GetExpenseByID(c.Request.Context(), uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.NotFoundResponse(c, "Expense not found")
			return
		}
		utils.InternalServerErrorResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Expense retrieved successfully", "data": expense})
}
