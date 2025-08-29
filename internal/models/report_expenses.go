package models

type ReportExpense struct {
	ReportID  uint `gorm:"primaryKey"`
	ExpenseID uint `gorm:"primaryKey"`
}
