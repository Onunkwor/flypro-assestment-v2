package services

import (
	"context"
	"errors"

	"github.com/onunkwor/flypro-assestment-v2/internal/models"
	"github.com/onunkwor/flypro-assestment-v2/internal/repository"
	"github.com/redis/go-redis/v9"
)

var (
	ErrInvalidReportState = errors.New("report cannot be submitted in current state")
	ErrInvalidOwnership   = errors.New("user does not have ownership")
)

type ReportService interface {
	CreateReport(ctx context.Context, report *models.ExpenseReport) error
	AddExpenseToReport(ctx context.Context, reportID, userId, expenseID uint) error
	SubmitReport(ctx context.Context, reportID, userID uint) error
	GetReportExpenses(ctx context.Context, userID uint, offset, limit int) ([]models.ExpenseReport, error)
}

type reportService struct {
	reportRepo  repository.ReportRepository
	expenseRepo repository.ExpenseRepository
	userRepo    repository.UserRepository
	redis       *redis.Client
}

func NewReportService(r repository.ReportRepository, exp repository.ExpenseRepository, userRepo repository.UserRepository, redis *redis.Client) *reportService {
	return &reportService{
		reportRepo:  r,
		expenseRepo: exp,
		userRepo:    userRepo,
		redis:       redis,
	}
}

func (s *reportService) CreateReport(ctx context.Context, report *models.ExpenseReport) error {
	_, err := s.userRepo.GetUserByID(ctx, report.UserID)
	if err != nil {
		return err
	}
	return s.reportRepo.CreateReport(ctx, report)
}

func (s *reportService) AddExpenseToReport(ctx context.Context, reportID, userID, expenseID uint) error {
	report, err := s.reportRepo.GetExpenseReportByID(ctx, reportID)
	if err != nil {
		return err
	}

	if report.UserID != userID {
		return ErrInvalidOwnership
	}

	expense, err := s.expenseRepo.GetExpenseByID(ctx, expenseID)
	if err != nil {
		return err
	}

	if expense.UserID != userID {
		return ErrInvalidOwnership
	}

	if err := s.reportRepo.AddExpenseToReport(ctx, reportID, expense); err != nil {
		return err
	}
	return s.reportRepo.IncrementReportTotal(ctx, reportID, expense.AmountUSD)

}

func (s *reportService) SubmitReport(ctx context.Context, reportID, userID uint) error {
	report, err := s.reportRepo.GetExpenseReportByID(ctx, reportID)
	if err != nil {
		return err
	}

	if report.UserID != userID {
		return ErrInvalidOwnership
	}

	if report.Status != "draft" {
		return ErrInvalidReportState
	}

	return s.reportRepo.SubmitReport(ctx, reportID)
}

func (s *reportService) GetReportExpenses(ctx context.Context, userID uint, offset, limit int) ([]models.ExpenseReport, error) {
	return s.reportRepo.GetReportExpenses(ctx, userID, offset, limit)
}
