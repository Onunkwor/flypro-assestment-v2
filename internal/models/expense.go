package models

type Expense struct {
	BaseModel
	UserID      uint    `json:"user_id" gorm:"not null"`
	Amount      float64 `json:"amount" gorm:"not null"`
	Currency    string  `json:"currency" gorm:"not null"`
	Category    string  `json:"category" gorm:"not null"`
	Description string  `json:"description"`
	Receipt     string  `json:"receipt"`
	Status      string  `json:"status" gorm:"default:'pending'"`
	User        *User   `json:"user" gorm:"foreignKey:UserID"`
}
