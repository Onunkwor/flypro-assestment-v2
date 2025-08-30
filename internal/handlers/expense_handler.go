package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/onunkwor/flypro-assestment-v2/internal/dto"
	"github.com/onunkwor/flypro-assestment-v2/internal/models"
	"github.com/onunkwor/flypro-assestment-v2/internal/repository"
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
		UserID:      request.UserId,
		Amount:      request.Amount,
		Currency:    request.Currency,
		Description: request.Description,
		Category:    request.Category,
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

func (h *ExpenseHandler) UpdateExpense(c *gin.Context) {
	var request dto.UpdateExpenseRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		formatted := utils.FormatValidationError(err)
		utils.ValidationErrorResponse(c, formatted)
		return
	}

	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		utils.BadRequestResponse(c, "Invalid expense ID")
		return
	}
	userIDParam := c.Query("user_id")
	if userIDParam == "" {
		utils.BadRequestResponse(c, "User ID is required")
		return
	}

	uid, err := strconv.ParseUint(userIDParam, 10, 64)
	if err != nil {
		utils.BadRequestResponse(c, "Invalid user ID")
		return
	}
	userID := uint(uid)
	expense := &models.Expense{
		Amount:      request.Amount,
		Currency:    request.Currency,
		Description: request.Description,
		Category:    request.Category,
	}

	if err := h.service.UpdateExpense(c.Request.Context(), uint(id), expense, userID); err != nil {
		if errors.Is(err, repository.ErrExpenseNotFound) {
			utils.NotFoundResponse(c, "Expense not found")
			return
		}
		utils.InternalServerErrorResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Expense updated successfully"})
}

func (h *ExpenseHandler) DeleteExpense(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		utils.BadRequestResponse(c, "Invalid expense ID")
		return
	}
	userIDParam := c.Query("user_id")
	if userIDParam == "" {
		utils.BadRequestResponse(c, "User ID is required")
		return
	}

	uid, err := strconv.ParseUint(userIDParam, 10, 64)
	if err != nil {
		utils.BadRequestResponse(c, "Invalid user ID")
		return
	}
	userID := uint(uid)
	if err := h.service.DeleteExpense(c.Request.Context(), uint(id), userID); err != nil {
		if errors.Is(err, repository.ErrExpenseNotFound) {
			utils.NotFoundResponse(c, "Expense not found")
			return
		}
		utils.InternalServerErrorResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Expense deleted successfully"})
}

func (h *ExpenseHandler) GetExpenses(c *gin.Context) {
	filters := make(map[string]interface{})
	if userIDParam := c.Query("user_id"); userIDParam != "" {
		userID, err := strconv.ParseUint(userIDParam, 10, 64)
		if err != nil {
			utils.BadRequestResponse(c, "Invalid user ID")
			return
		}
		filters["user_id"] = uint(userID)
	}

	if category := c.Query("category"); category != "" {
		filters["category"] = utils.NormalizeCategory(category)
	}

	var (
		limit  = 20
		offset = 0
		err    error
	)

	if limitStr := c.Query("limit"); limitStr != "" {
		limit, err = strconv.Atoi(limitStr)
		if err != nil || limit <= 0 {
			utils.BadRequestResponse(c, "Invalid limit value")
			return
		}
	}

	if offsetStr := c.Query("offset"); offsetStr != "" {
		offset, err = strconv.Atoi(offsetStr)
		if err != nil || offset < 0 {
			utils.BadRequestResponse(c, "Invalid offset value")
			return
		}
	}

	expenses, err := h.service.GetExpenses(c.Request.Context(), filters, offset, limit)
	if err != nil {
		utils.InternalServerErrorResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":   expenses,
		"count":  len(expenses),
		"offset": offset,
		"limit":  limit,
	})
}
