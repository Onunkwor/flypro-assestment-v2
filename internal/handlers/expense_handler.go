package handlers

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/onunkwor/flypro-assestment-v2/internal/dto"
	"github.com/onunkwor/flypro-assestment-v2/internal/utils"
)

type ExpenseHandler struct {
	service interface{}
}

func NewExpenseHandler(service interface{}) *ExpenseHandler {
	return &ExpenseHandler{service: service}
}

func (h *ExpenseHandler) CreateExpense(c *gin.Context) {
	var request dto.CreateExpenseRequest
	request.Currency = strings.ToUpper(request.Currency)
	if err := c.ShouldBindJSON(&request); err != nil {
		formatted := utils.FormatValidationError(err)
		utils.ValidationErrorResponse(c, formatted)
		return
	}
}
