package repository

import (
	"context"
	"errors"

	"github.com/onunkwor/flypro-assestment-v2/internal/models"
	"gorm.io/gorm"
)

var ErrReportNotFound = errors.New("report not found")

type ReportRepository interface {
	CreateReport(ctx context.Context, report *models.ExpenseReport) error
	AddExpenseToReport(ctx context.Context, reportID uint, expense *models.Expense) error
	GetExpenseReportByID(ctx context.Context, id uint) (*models.ExpenseReport, error)
	IncrementReportTotal(ctx context.Context, reportID uint, amount float64) error
	GetReportExpenses(ctx context.Context, userID uint, offset, limit int) ([]models.ExpenseReport, error)
	SubmitReport(ctx context.Context, reportID uint) error
}

type reportRepo struct {
	db *gorm.DB
}

func NewReportRepository(db *gorm.DB) ReportRepository {
	return &reportRepo{db: db}
}

func (r *reportRepo) CreateReport(ctx context.Context, report *models.ExpenseReport) error {
	return r.db.WithContext(ctx).Create(report).Error
}

func (r *reportRepo) AddExpenseToReport(ctx context.Context, reportID uint, expense *models.Expense) error {

	return r.db.WithContext(ctx).
		Model(&models.ExpenseReport{BaseModel: models.BaseModel{ID: reportID}}).
		Association("Expenses").
		Append(expense)
}

func (r *reportRepo) GetExpenseReportByID(ctx context.Context, id uint) (*models.ExpenseReport, error) {
	var report models.ExpenseReport
	if err := r.db.WithContext(ctx).Preload("Expenses").First(&report, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrReportNotFound
		}
		return nil, err
	}
	return &report, nil
}

func (r *reportRepo) IncrementReportTotal(ctx context.Context, reportID uint, amount float64) error {
	return r.db.WithContext(ctx).
		Model(&models.ExpenseReport{}).
		Where("id = ?", reportID).
		UpdateColumn("total", gorm.Expr("total + ?", amount)).
		Error
}

func (r *reportRepo) GetReportExpenses(ctx context.Context, userID uint, offset, limit int) ([]models.ExpenseReport, error) {
	var reports []models.ExpenseReport
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Offset(offset).
		Limit(limit).
		Preload("Expenses").
		Preload("User").
		Find(&reports).Error
	return reports, err
}

func (r *reportRepo) SubmitReport(ctx context.Context, reportID uint) error {
	return r.db.WithContext(ctx).
		Model(&models.ExpenseReport{}).
		Where("id = ?", reportID).
		UpdateColumn("status", "submitted").
		Error
}
