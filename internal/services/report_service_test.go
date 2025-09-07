package services_test

import (
	"context"
	"errors"
	"testing"

	"github.com/onunkwor/flypro-assestment-v2/internal/models"
	"github.com/onunkwor/flypro-assestment-v2/internal/repository"
	"github.com/onunkwor/flypro-assestment-v2/internal/services"
	"github.com/onunkwor/flypro-assestment-v2/tests/mocks"
	"go.uber.org/mock/gomock"
)

func TestCreateReport_InvalidUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockReportRepo := mocks.NewMockReportRepository(ctrl)
	mockUserRepo := mocks.NewMockUserRepository(ctrl)

	service := services.NewReportService(mockReportRepo, nil, mockUserRepo, nil)

	report := &models.ExpenseReport{
		UserID: 1,
		Title:  "Test Report",
	}

	mockUserRepo.EXPECT().GetUserByID(gomock.Any(), report.UserID).Return(nil, repository.ErrUserNotFound)

	mockReportRepo.EXPECT().
		CreateReport(gomock.Any(), gomock.Any()).
		Times(0)

	err := service.CreateReport(context.Background(), report)

	if !errors.Is(err, repository.ErrUserNotFound) {
		t.Fatalf("expected ErrUserNotFound, got %v", err)
	}
}

func TestCreateReport_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockReportRepo := mocks.NewMockReportRepository(ctrl)
	mockUserRepo := mocks.NewMockUserRepository(ctrl)

	service := services.NewReportService(mockReportRepo, nil, mockUserRepo, nil)

	report := &models.ExpenseReport{
		UserID: 1,
		Title:  "Valid Report",
	}

	mockUser := &models.User{BaseModel: models.BaseModel{ID: report.UserID}, Name: "Test User"}

	mockUserRepo.EXPECT().
		GetUserByID(gomock.Any(), report.UserID).
		Return(mockUser, nil)

	mockReportRepo.EXPECT().
		CreateReport(gomock.Any(), report).
		Return(nil).
		Times(1)

	err := service.CreateReport(context.Background(), report)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestAddExpenseToReport_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockReportRepo := mocks.NewMockReportRepository(ctrl)
	mockExpenseRepo := mocks.NewMockExpenseRepository(ctrl)
	mockUserRepo := mocks.NewMockUserRepository(ctrl)

	service := services.NewReportService(mockReportRepo, mockExpenseRepo, mockUserRepo, nil)

	reportID := uint(1)
	userID := uint(1)
	expenseID := uint(10)

	report := &models.ExpenseReport{BaseModel: models.BaseModel{ID: reportID}, UserID: userID, Total: 100}
	expense := &models.Expense{BaseModel: models.BaseModel{ID: expenseID}, UserID: userID, AmountUSD: 50}

	mockReportRepo.EXPECT().GetExpenseReportByID(gomock.Any(), reportID).Return(report, nil)
	mockExpenseRepo.EXPECT().GetExpenseByID(gomock.Any(), expenseID).Return(expense, nil)
	mockReportRepo.EXPECT().AddExpenseToReportWithTotal(gomock.Any(), reportID, expense).Return(nil)

	err := service.AddExpenseToReport(context.Background(), expenseID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestSubmitReport_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockReportRepo := mocks.NewMockReportRepository(ctrl)
	service := services.NewReportService(mockReportRepo, nil, nil, nil)

	reportID := uint(1)
	userID := uint(1)
	report := &models.ExpenseReport{BaseModel: models.BaseModel{ID: reportID}, UserID: userID, Status: "draft"}

	mockReportRepo.EXPECT().GetExpenseReportByID(gomock.Any(), reportID).Return(report, nil)
	mockReportRepo.EXPECT().SubmitReport(gomock.Any(), reportID).Return(nil)

	err := service.SubmitReport(context.Background(), reportID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestGetReportExpenses_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockReportRepo := mocks.NewMockReportRepository(ctrl)
	service := services.NewReportService(mockReportRepo, nil, nil, nil)

	userID := uint(1)
	offset := 0
	limit := 10
	reports := []models.ExpenseReport{
		{BaseModel: models.BaseModel{ID: 1}, UserID: userID, Title: "Report 1"},
		{BaseModel: models.BaseModel{ID: 2}, UserID: userID, Title: "Report 2"},
	}

	mockReportRepo.EXPECT().GetReportExpenses(gomock.Any(), userID, offset, limit).Return(reports, nil)

	result, err := service.GetReportExpenses(context.Background(), userID, offset, limit)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(result) != len(reports) {
		t.Fatalf("expected %d reports, got %d", len(reports), len(result))
	}
}
