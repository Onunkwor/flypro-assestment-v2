package services_test

import (
	"context"
	"testing"

	"github.com/onunkwor/flypro-assestment-v2/internal/models"
	"github.com/onunkwor/flypro-assestment-v2/internal/services"
	"github.com/onunkwor/flypro-assestment-v2/tests/mocks"
	"go.uber.org/mock/gomock"
)

func TestCreateExpanse_CurrencyConversionFAil(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCurrency := mocks.NewMockCurrencyConverter(ctrl)
	mockRepo := mocks.NewMockExpenseRepository(ctrl)
	srv := services.NewExpenseService(nil, mockCurrency, mockRepo)

	exp := &models.Expense{
		UserID:   1,
		Amount:   100,
		Currency: "NGN",
	}

	mockCurrency.EXPECT().
		Convert(gomock.Any(), 100.0, "NGN", "USD").
		Return(0.0, 0.0, services.ErrCurrencyConversionFailed)

	err := srv.CreateExpense(context.Background(), exp)
	if err != services.ErrCurrencyConversionFailed {
		t.Errorf("expected error %v, got %v", services.ErrCurrencyConversionFailed, err)
	}
}

func TestCreateExpense_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockExpenseRepository(ctrl)
	mockCurrency := mocks.NewMockCurrencyConverter(ctrl)

	srv := services.NewExpenseService(nil, mockCurrency, mockRepo)

	exp := &models.Expense{
		UserID:   1,
		Amount:   100,
		Currency: "NGN",
	}

	mockCurrency.EXPECT().
		Convert(gomock.Any(), exp.Amount, "NGN", "USD").
		Return(0.26, 0.0026, nil)

	mockRepo.EXPECT().
		Create(gomock.Any(), gomock.AssignableToTypeOf(&models.Expense{})).
		DoAndReturn(func(ctx context.Context, e *models.Expense) error {
			if e.AmountUSD != 0.26 {
				t.Fatalf("expected AmountUSD 0.26, got %v", e.AmountUSD)
			}
			if e.ExchangeRate != 0.0026 {
				t.Fatalf("expected ExchangeRate 0.0026, got %v", e.ExchangeRate)
			}
			return nil
		})

	if err := srv.CreateExpense(context.Background(), exp); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
