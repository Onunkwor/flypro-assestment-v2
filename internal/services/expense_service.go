package services

import (
	"context"
	"errors"
	"strings"

	"github.com/onunkwor/flypro-assestment-v2/internal/models"
	"github.com/onunkwor/flypro-assestment-v2/internal/repository"
	"github.com/redis/go-redis/v9"
)

var ErrCurrencyConversionFailed = errors.New("currency conversion failed")

type ExpenseService struct {
	repo        repository.ExpenseRepository
	redis       *redis.Client
	currencySvc CurrencyConverter
}

func NewExpenseService(redis *redis.Client, currencySvc CurrencyConverter, repo repository.ExpenseRepository) *ExpenseService {
	return &ExpenseService{repo: repo, redis: redis, currencySvc: currencySvc}
}

func (s *ExpenseService) CreateExpense(ctx context.Context, expense *models.Expense) error {
	currency := strings.ToUpper(expense.Currency)

	if currency == "USD" {

		expense.AmountUSD = expense.Amount
		expense.ExchangeRate = 1.0
	} else {
		convertedAmount, rate, err := s.currencySvc.Convert(ctx, expense.Amount, currency, "USD")
		if err != nil {
			return ErrCurrencyConversionFailed
		}
		expense.AmountUSD = convertedAmount
		expense.ExchangeRate = rate
	}

	return s.repo.Create(ctx, expense)
}

func (s *ExpenseService) GetExpenseByID(ctx context.Context, id uint) (*models.Expense, error) {
	return s.repo.GetExpenseByID(ctx, id)
}

func (s *ExpenseService) UpdateExpense(ctx context.Context, id uint, expense *models.Expense) error {
	currency := strings.ToUpper(expense.Currency)

	if currency == "USD" {
		expense.AmountUSD = expense.Amount
		expense.ExchangeRate = 1.0
	} else {
		convertedAmount, rate, err := s.currencySvc.Convert(ctx, expense.Amount, currency, "USD")
		if err != nil {
			return ErrCurrencyConversionFailed
		}
		expense.AmountUSD = convertedAmount
		expense.ExchangeRate = rate
	}

	return s.repo.UpdateExpense(ctx, id, expense)
}

func (s *ExpenseService) DeleteExpense(ctx context.Context, id uint) error {
	return s.repo.DeleteExpense(ctx, id)
}

func (s *ExpenseService) GetExpenses(ctx context.Context, filters map[string]interface{}, offset, limit int) ([]models.Expense, error) {
	return s.repo.GetExpenses(ctx, filters, offset, limit)
}
