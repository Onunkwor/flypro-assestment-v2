package services

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/onunkwor/flypro-assestment-v2/internal/models"
	"github.com/onunkwor/flypro-assestment-v2/internal/repository"
	"github.com/onunkwor/flypro-assestment-v2/internal/utils"
	"github.com/redis/go-redis/v9"
)

var ErrCurrencyConversionFailed = errors.New("currency conversion failed")

type ExpenseService interface {
	CreateExpense(ctx context.Context, expense *models.Expense) error
	GetExpenseByID(ctx context.Context, id uint) (*models.Expense, error)
	UpdateExpense(ctx context.Context, id uint, expense *models.Expense, userId uint) error
	DeleteExpense(ctx context.Context, id uint, userId uint) error
	GetExpenses(ctx context.Context, filters map[string]interface{}, offset, limit int) ([]models.Expense, error)
}

type expenseSrv struct {
	repo        repository.ExpenseRepository
	redis       RedisClient
	currencySvc CurrencyConverter
}

func NewExpenseService(redis RedisClient, currencySvc CurrencyConverter, repo repository.ExpenseRepository) ExpenseService {
	return &expenseSrv{repo: repo, redis: redis, currencySvc: currencySvc}
}

func (s *expenseSrv) CreateExpense(ctx context.Context, expense *models.Expense) error {
	currency := strings.ToUpper(expense.Currency)

	if expense.Currency == "USD" {

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
	s.invalidateExpensesCache(ctx)
	return s.repo.Create(ctx, expense)
}

func (s *expenseSrv) GetExpenseByID(ctx context.Context, id uint) (*models.Expense, error) {
	return s.repo.GetExpenseByID(ctx, id)
}

func (s *expenseSrv) UpdateExpense(ctx context.Context, id uint, expense *models.Expense, userId uint) error {
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
	s.invalidateExpensesCache(ctx)
	return s.repo.UpdateExpense(ctx, id, expense, userId)
}

func (s *expenseSrv) DeleteExpense(ctx context.Context, id uint, userId uint) error {
	s.invalidateExpensesCache(ctx)
	return s.repo.DeleteExpense(ctx, id, userId)
}

func (s *expenseSrv) GetExpenses(ctx context.Context, filters map[string]interface{}, offset, limit int) ([]models.Expense, error) {
	key := utils.ExpensesCacheKey(filters, offset, limit)
	if s.redis != nil {
		val, err := s.redis.Get(ctx, key).Result()
		if err == nil {
			var expenses []models.Expense
			if unmarshalErr := json.Unmarshal([]byte(val), &expenses); unmarshalErr == nil {
				return expenses, nil
			}
			_ = s.redis.Del(ctx, key).Err()
		} else if err != redis.Nil {
			return nil, err
		}
	}
	expenses, err := s.repo.GetExpenses(ctx, filters, offset, limit)
	if err != nil {
		return nil, err
	}
	if s.redis != nil {
		bytes, _ := json.Marshal(expenses)
		_ = s.redis.Set(ctx, key, bytes, time.Minute*30).Err()
	}

	return expenses, nil
}

func (s *expenseSrv) invalidateExpensesCache(ctx context.Context) {
	if s.redis == nil {
		return
	}
	iter := s.redis.Scan(ctx, 0, "expenses:*", 0).Iterator()
	for iter.Next(ctx) {
		_ = s.redis.Del(ctx, iter.Val()).Err()
	}
}
