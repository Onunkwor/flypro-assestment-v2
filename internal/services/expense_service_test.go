package services_test

import (
	"context"
	"testing"

	"github.com/onunkwor/flypro-assestment-v2/internal/models"
	"github.com/onunkwor/flypro-assestment-v2/internal/services"
	"github.com/onunkwor/flypro-assestment-v2/tests/mocks"
	"go.uber.org/mock/gomock"
)

func TestCreateExpense(t *testing.T) {
	tests := []struct {
		name         string
		expense      *models.Expense
		mockRepo     func(repo *mocks.MockExpenseRepository)
		mockCurrency func(ctrl *gomock.Controller) *mocks.MockCurrencyConverter
		expectedErr  error
		assert       func(t *testing.T, exp *models.Expense)
	}{
		{
			name: "USD_Currency",
			expense: &models.Expense{
				Currency: "USD",
				Amount:   100,
			},
			mockRepo: func(repo *mocks.MockExpenseRepository) {
				repo.EXPECT().
					Create(gomock.Any(), gomock.Any()).
					Return(nil)
			},
			mockCurrency: func(ctrl *gomock.Controller) *mocks.MockCurrencyConverter {
				return mocks.NewMockCurrencyConverter(ctrl)
			},
			expectedErr: nil,
			assert: func(t *testing.T, exp *models.Expense) {
				if exp.AmountUSD != 100 || exp.ExchangeRate != 1.0 {
					t.Errorf("expected AmountUSD 100 and ExchangeRate 1.0, got %v and %v", exp.AmountUSD, exp.ExchangeRate)
				}
			},
		},
		{
			name: "NonUSD_Currency_ConversionSuccess",
			expense: &models.Expense{
				Currency: "EUR",
				Amount:   200,
			},
			mockRepo: func(repo *mocks.MockExpenseRepository) {
				repo.EXPECT().
					Create(gomock.Any(), gomock.Any()).
					Return(nil)
			},
			mockCurrency: func(ctrl *gomock.Controller) *mocks.MockCurrencyConverter {
				mockCurr := mocks.NewMockCurrencyConverter(ctrl)
				mockCurr.EXPECT().
					Convert(gomock.Any(), 200.0, "EUR", "USD").
					Return(220.0, 1.1, nil)
				return mockCurr
			},
			expectedErr: nil,
			assert: func(t *testing.T, exp *models.Expense) {
				if exp.AmountUSD != 220.0 || exp.ExchangeRate != 1.1 {
					t.Errorf("expected AmountUSD=220, ExchangeRate=1.1, got %+v", exp)
				}
			},
		},
		{
			name: "NonUSD_Currency_ConversionFails",
			expense: &models.Expense{
				Currency: "EUR",
				Amount:   200,
			},
			mockRepo: func(repo *mocks.MockExpenseRepository) {

			},
			mockCurrency: func(ctrl *gomock.Controller) *mocks.MockCurrencyConverter {
				mockCurr := mocks.NewMockCurrencyConverter(ctrl)
				mockCurr.EXPECT().
					Convert(gomock.Any(), float64(200), "EUR", "USD").
					Return(0.0, 0.0, services.ErrCurrencyConversionFailed)
				return mockCurr
			},
			expectedErr: services.ErrCurrencyConversionFailed,
			assert: func(t *testing.T, exp *models.Expense) {
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := mocks.NewMockExpenseRepository(ctrl)
			mockCurr := tt.mockCurrency(ctrl)

			tt.mockRepo(mockRepo)

			svc := services.NewExpenseService(nil, mockCurr, mockRepo)

			err := svc.CreateExpense(context.Background(), tt.expense)
			if (tt.expectedErr != nil && (err == nil || err.Error() != tt.expectedErr.Error())) ||
				(tt.expectedErr == nil && err != nil) {
				t.Errorf("expected error %v, got %v", tt.expectedErr, err)
			}

			tt.assert(t, tt.expense)
		})
	}
}
