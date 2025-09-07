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
	AddExpenseToReport(ctx context.Context, reportID, expenseID uint) error
	SubmitReport(ctx context.Context, reportID uint) error
	GetReportExpenses(ctx context.Context, userID uint, offset, limit int) ([]models.ExpenseReport, error)
}

type reportService struct {
	reportRepo  repository.ReportRepository
	expenseRepo repository.ExpenseRepository
	userRepo    repository.UserRepository
	redis       *redis.Client
}

func NewReportService(r repository.ReportRepository, e repository.ExpenseRepository, u repository.UserRepository, redis *redis.Client) *reportService {
	return &reportService{
		reportRepo:  r,
		expenseRepo: e,
		userRepo:    u,
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

func (s *reportService) AddExpenseToReport(ctx context.Context, reportID uint, expenseID uint) error {
	expense, err := s.expenseRepo.GetExpenseByID(ctx, expenseID)
	if err != nil {
		return err
	}

	return s.reportRepo.AddExpenseToReportWithTotal(ctx, reportID, expense)
}

func (s *reportService) SubmitReport(ctx context.Context, reportID uint) error {
	report, err := s.reportRepo.GetExpenseReportByID(ctx, reportID)
	if err != nil {
		return err
	}
	if report.Status != "draft" {
		return ErrInvalidReportState
	}
	return s.reportRepo.SubmitReport(ctx, reportID)
}

func (s *reportService) GetReportExpenses(ctx context.Context, userID uint, offset, limit int) ([]models.ExpenseReport, error) {
	return s.reportRepo.GetReportExpenses(ctx, userID, offset, limit)
}
