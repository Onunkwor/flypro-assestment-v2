package dto

type CreateExpenseRequest struct {
	UserId      uint    `json:"user_id" binding:"required"`
	Amount      float64 `json:"amount" binding:"required,gt=0"`
	Currency    string  `json:"currency" binding:"required,len=3,oneof=USD EUR GBP NGN"`
	Category    string  `json:"category" binding:"required,oneof=travel meals office supplies"`
	Description string  `json:"description" binding:"max=500"`
}

type UpdateExpenseRequest struct {
	UserId      uint    `json:"user_id" binding:"required"`
	Amount      float64 `json:"amount" binding:"required,gt=0"`
	Currency    string  `json:"currency" binding:"required,len=3,oneof=USD EUR GBP NGN"`
	Category    string  `json:"category" binding:"required,oneof=travel meals office supplies"`
	Description string  `json:"description" binding:"max=500"`
}
