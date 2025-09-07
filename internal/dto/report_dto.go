package dto

import "github.com/onunkwor/flypro-assestment-v2/internal/utils"

type CreateReportRequest struct {
	Title  string `json:"title" binding:"required"`
	UserID uint   `json:"user_id" binding:"required"`
}

type AddExpenseToReportRequest struct {
	ExpenseID uint `json:"expense_id" binding:"required"`
}

func (r *CreateReportRequest) Sanitize() {
	r.Title = utils.SanitizeString(r.Title)
}
