package models

type User struct {
	BaseModel
	Email string `json:"email" gorm:"uniqueIndex;not null"`
	Name  string `json:"name" gorm:"not null"`
}
