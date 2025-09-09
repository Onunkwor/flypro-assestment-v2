package services_test

import (
	"context"
	"errors"
	"testing"

	"github.com/onunkwor/flypro-assestment-v2/internal/models"
	"github.com/onunkwor/flypro-assestment-v2/internal/services"
	"github.com/onunkwor/flypro-assestment-v2/tests/mocks"
	"go.uber.org/mock/gomock"
)

func TestCreateReport(t *testing.T) {
	tests := []struct {
		name        string
		report      *models.ExpenseReport
		mockUser    func(repo *mocks.MockUserRepository)
		mockReport  func(repo *mocks.MockReportRepository)
		expectedErr error
	}{
		{
			name:   "Success",
			report: &models.ExpenseReport{UserID: 1, Title: "Trip"},
			mockUser: func(repo *mocks.MockUserRepository) {
				repo.EXPECT().GetUserByID(gomock.Any(), uint(1)).Return(&models.User{BaseModel: models.BaseModel{ID: 1}}, nil)
			},
			mockReport: func(repo *mocks.MockReportRepository) {
				repo.EXPECT().CreateReport(gomock.Any(), gomock.Any()).Return(nil)
			},
			expectedErr: nil,
		},
		{
			name:   "UserNotFound",
			report: &models.ExpenseReport{UserID: 2, Title: "Trip"},
			mockUser: func(repo *mocks.MockUserRepository) {
				repo.EXPECT().GetUserByID(gomock.Any(), uint(2)).Return(nil, errors.New("not found"))
			},
			mockReport:  func(repo *mocks.MockReportRepository) {},
			expectedErr: errors.New("not found"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockReportRepo := mocks.NewMockReportRepository(ctrl)
			mockUserRepo := mocks.NewMockUserRepository(ctrl)

			tt.mockUser(mockUserRepo)
			tt.mockReport(mockReportRepo)

			service := services.NewReportService(mockReportRepo, nil, mockUserRepo, nil)

			err := service.CreateReport(context.Background(), tt.report)
			if (tt.expectedErr != nil && (err == nil || err.Error() != tt.expectedErr.Error())) ||
				(tt.expectedErr == nil && err != nil) {
				t.Fatalf("expected error %v, got %v", tt.expectedErr, err)
			}
		})
	}
}

func TestAddExpenseToReport(t *testing.T) {
	tests := []struct {
		name        string
		reportID    uint
		expense     *models.Expense
		mockReport  func(repo *mocks.MockReportRepository)
		expectedErr error
	}{
		{
			name:     "Success",
			reportID: 1,
			expense:  &models.Expense{BaseModel: models.BaseModel{ID: 1}, UserID: 1, AmountUSD: 100},
			mockReport: func(repo *mocks.MockReportRepository) {
				repo.EXPECT().AddExpenseToReportWithTotal(gomock.Any(), uint(1), gomock.Any()).Return(nil)
			},
			expectedErr: nil,
		},
		{
			name:     "RepoFailure",
			reportID: 1,
			expense:  &models.Expense{BaseModel: models.BaseModel{ID: 2}, UserID: 2, AmountUSD: 50},
			mockReport: func(repo *mocks.MockReportRepository) {
				repo.EXPECT().AddExpenseToReportWithTotal(gomock.Any(), uint(1), gomock.Any()).Return(errors.New("db error"))
			},
			expectedErr: errors.New("db error"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockReportRepo := mocks.NewMockReportRepository(ctrl)
			service := services.NewReportService(mockReportRepo, nil, nil, nil)

			tt.mockReport(mockReportRepo)

			err := service.AddExpenseToReport(context.Background(), tt.reportID, tt.expense)
			if (tt.expectedErr != nil && (err == nil || err.Error() != tt.expectedErr.Error())) ||
				(tt.expectedErr == nil && err != nil) {
				t.Fatalf("expected %v, got %v", tt.expectedErr, err)
			}
		})
	}
}

func TestSubmitReport(t *testing.T) {
	tests := []struct {
		name        string
		reportID    uint
		mockReport  func(repo *mocks.MockReportRepository)
		expectedErr error
	}{
		{
			name:     "Success",
			reportID: 1,
			mockReport: func(repo *mocks.MockReportRepository) {
				report := &models.ExpenseReport{BaseModel: models.BaseModel{ID: 1}, Status: "draft"}
				repo.EXPECT().GetExpenseReportByID(gomock.Any(), uint(1)).Return(report, nil)
				repo.EXPECT().SubmitReport(gomock.Any(), uint(1)).Return(nil)
			},
			expectedErr: nil,
		},
		{
			name:     "ReportNotFound",
			reportID: 2,
			mockReport: func(repo *mocks.MockReportRepository) {
				repo.EXPECT().GetExpenseReportByID(gomock.Any(), uint(2)).Return(nil, errors.New("not found"))
			},
			expectedErr: errors.New("not found"),
		},
		{
			name:     "InvalidStatus",
			reportID: 3,
			mockReport: func(repo *mocks.MockReportRepository) {
				report := &models.ExpenseReport{BaseModel: models.BaseModel{ID: 3}, Status: "submitted"}
				repo.EXPECT().GetExpenseReportByID(gomock.Any(), uint(3)).Return(report, nil)
			},
			expectedErr: services.ErrInvalidReportState,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockReportRepo := mocks.NewMockReportRepository(ctrl)
			service := services.NewReportService(mockReportRepo, nil, nil, nil)

			tt.mockReport(mockReportRepo)

			err := service.SubmitReport(context.Background(), tt.reportID)
			if (tt.expectedErr != nil && (err == nil || err.Error() != tt.expectedErr.Error())) ||
				(tt.expectedErr == nil && err != nil) {
				t.Fatalf("expected %v, got %v", tt.expectedErr, err)
			}
		})
	}
}

func TestGetReportExpenses(t *testing.T) {
	tests := []struct {
		name        string
		userID      uint
		offset      int
		limit       int
		mockReport  func(repo *mocks.MockReportRepository)
		expected    []models.ExpenseReport
		expectedErr error
	}{
		{
			name:   "Success",
			userID: 1, offset: 0, limit: 10,
			mockReport: func(repo *mocks.MockReportRepository) {
				reports := []models.ExpenseReport{
					{BaseModel: models.BaseModel{ID: 1}, UserID: 1, Title: "Report 1"},
					{BaseModel: models.BaseModel{ID: 2}, UserID: 1, Title: "Report 2"},
				}
				repo.EXPECT().GetReportExpenses(gomock.Any(), uint(1), 0, 10).Return(reports, nil)
			},
			expected: []models.ExpenseReport{
				{BaseModel: models.BaseModel{ID: 1}, UserID: 1, Title: "Report 1"},
				{BaseModel: models.BaseModel{ID: 2}, UserID: 1, Title: "Report 2"},
			},
			expectedErr: nil,
		},
		{
			name:   "RepoError",
			userID: 2, offset: 0, limit: 5,
			mockReport: func(repo *mocks.MockReportRepository) {
				repo.EXPECT().GetReportExpenses(gomock.Any(), uint(2), 0, 5).Return(nil, errors.New("db error"))
			},
			expected:    nil,
			expectedErr: errors.New("db error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockReportRepo := mocks.NewMockReportRepository(ctrl)
			service := services.NewReportService(mockReportRepo, nil, nil, nil)

			tt.mockReport(mockReportRepo)

			result, err := service.GetReportExpenses(context.Background(), tt.userID, tt.offset, tt.limit)
			if (tt.expectedErr != nil && (err == nil || err.Error() != tt.expectedErr.Error())) ||
				(tt.expectedErr == nil && err != nil) {
				t.Fatalf("expected error %v, got %v", tt.expectedErr, err)
			}

			if len(result) != len(tt.expected) {
				t.Fatalf("expected %d reports, got %d", len(tt.expected), len(result))
			}
		})
	}
}
