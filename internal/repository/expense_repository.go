package repository

import (
	"context"
	"errors"

	"github.com/onunkwor/flypro-assestment-v2/internal/models"
	"gorm.io/gorm"
)

var ErrExpenseNotFound = errors.New("expense not found")

type ExpenseRepository interface {
	Create(ctx context.Context, expense *models.Expense) error
	GetExpenseByID(ctx context.Context, id uint) (*models.Expense, error)
	GetExpenses(ctx context.Context, filters map[string]interface{}, offset, limit int) ([]models.Expense, error)
	UpdateExpense(ctx context.Context, id uint, expense *models.Expense) error
	DeleteExpense(ctx context.Context, id uint) error
}

type expenseRepo struct {
	db *gorm.DB
}

func NewExpenseRepository(db *gorm.DB) ExpenseRepository {
	return &expenseRepo{db: db}
}

func (r *expenseRepo) Create(ctx context.Context, expense *models.Expense) error {

	return r.db.WithContext(ctx).Create(expense).Error
}

func (r *expenseRepo) GetExpenseByID(ctx context.Context, id uint) (*models.Expense, error) {
	var expense models.Expense
	if err := r.db.WithContext(ctx).Preload("User").First(&expense, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrExpenseNotFound
		}
		return nil, err
	}
	return &expense, nil
}

func (r *expenseRepo) GetExpenses(ctx context.Context, filters map[string]interface{}, offset, limit int) ([]models.Expense, error) {
	var expenses []models.Expense
	query := r.db.WithContext(ctx).Model(&models.Expense{}).Preload("User")

	allowedFilters := map[string]bool{"user_id": true, "category": true, "status": true}
	for key, value := range filters {
		if allowedFilters[key] {
			query = query.Where(key+" = ?", value)
		}
	}

	if err := query.Offset(offset).Limit(limit).Find(&expenses).Error; err != nil {
		return nil, err
	}
	return expenses, nil
}

func (r *expenseRepo) UpdateExpense(ctx context.Context, id uint, expense *models.Expense) error {
	result := r.db.WithContext(ctx).Model(&models.Expense{}).Where("id = ?", id).Updates(expense)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrExpenseNotFound
	}
	return nil
}

func (r *expenseRepo) DeleteExpense(ctx context.Context, id uint) error {
	result := r.db.WithContext(ctx).Delete(&models.Expense{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrExpenseNotFound
	}
	return nil
}
