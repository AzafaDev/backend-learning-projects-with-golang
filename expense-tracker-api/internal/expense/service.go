package expense

import (
	"context"
	"expense-tracker-api/internal/httpserver"
	"fmt"
	"net/http"
	"time"

	"github.com/go-playground/validator"
	"github.com/google/uuid"
)

const dateLayout = "2006-01-02"

type ExpenseRepo interface {
	GetExpenseByID(ctx context.Context, id, userID uuid.UUID) (*Expense, error)
	GetExpenses(ctx context.Context, userID uuid.UUID, filter ExpenseFilter, pagination Pagination) ([]Expense, int, error)
	CreateExpense(ctx context.Context, userID uuid.UUID, req CreateExpenseRequest) (*Expense, error)
	UpdateExpense(ctx context.Context, id, userID uuid.UUID, req UpdateExpenseRequest) (*Expense, error)
	DeleteExpense(ctx context.Context, id, userID uuid.UUID) error
}

type ListExpensesQuery struct {
	Period    string
	StartDate string
	EndDate   string
	Category  string
	Page      int
	Limit     int
}

type ExpenseService struct {
	Repo     ExpenseRepo
	Validate *validator.Validate
}

func NewExpenseService(repo ExpenseRepo) *ExpenseService {
	return &ExpenseService{
		Repo:     repo,
		Validate: validator.New(),
	}
}

func (e *ExpenseService) CreateExpense(ctx context.Context, userID uuid.UUID, req CreateExpenseRequest) (*Expense, error) {
	if err := e.Validate.Struct(req); err != nil {
		return nil, httpserver.NewClientError(http.StatusBadRequest, formatValidationError(err))
	}

	return e.Repo.CreateExpense(ctx, userID, req)
}

func (e *ExpenseService) GetExpenseByID(ctx context.Context, userID, id uuid.UUID) (*Expense, error) {
	return e.Repo.GetExpenseByID(ctx, id, userID)
}

func (e *ExpenseService) UpdateExpense(ctx context.Context, userID, id uuid.UUID, req UpdateExpenseRequest) (*Expense, error) {
	if err := e.Validate.Struct(req); err != nil {
		return nil, httpserver.NewClientError(http.StatusBadRequest, formatValidationError(err))
	}

	return e.Repo.UpdateExpense(ctx, id, userID, req)
}

func (e *ExpenseService) DeleteExpense(ctx context.Context, userID, id uuid.UUID) error {
	return e.Repo.DeleteExpense(ctx, id, userID)
}

func (e *ExpenseService) GetExpenses(ctx context.Context, userID uuid.UUID, q ListExpensesQuery) (*ExpenseListResponse, error) {
	start, end, err := calculatePeriodRange(q.Period, q.StartDate, q.EndDate, time.Now().UTC())
	if err != nil {
		return nil, httpserver.NewClientError(http.StatusBadRequest, err.Error())
	}

	filter := ExpenseFilter{StartDate: start, EndDate: end}
	if q.Category != "" {
		category := ExpenseCategory(q.Category)
		if err := e.Validate.Var(category, "oneof=Groceries Leisure Electronics Utilities Clothing Health Others"); err != nil {
			return nil, httpserver.NewClientError(http.StatusBadRequest, "category must be one of 'Groceries','Leisure','Electronics','Utilities','Clothing','Health','Others'")
		}
		filter.Category = &category
	}

	pagination := Pagination{Page: q.Page, Limit: q.Limit}
	if pagination.Page < 1 {
		pagination.Page = 1
	}
	if pagination.Limit < 1 {
		pagination.Limit = 20
	}
	if pagination.Limit > 100 {
		pagination.Limit = 100
	}

	expenses, total, err := e.Repo.GetExpenses(ctx, userID, filter, pagination)
	if err != nil {
		return nil, err
	}
	if expenses == nil {
		expenses = []Expense{}
	}

	return &ExpenseListResponse{
		Data:  expenses,
		Page:  pagination.Page,
		Limit: pagination.Limit,
		Total: total,
	}, nil
}

func calculatePeriodRange(period, startDateStr, endDateStr string, now time.Time) (start, end *time.Time, err error) {
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)

	switch period {
	case "":
		return nil, nil, nil
	case "week":
		s := today.AddDate(0, 0, -7)
		return &s, &today, nil
	case "month":
		s := today.AddDate(0, -1, 0)
		return &s, &today, nil
	case "3months":
		s := today.AddDate(0, -3, 0)
		return &s, &today, nil
	case "custom":
		if startDateStr == "" || endDateStr == "" {
			return nil, nil, fmt.Errorf("start_date and end_date are required when period=custom")
		}
		parsedStart, err := time.Parse(dateLayout, startDateStr)
		if err != nil {
			return nil, nil, fmt.Errorf("start_date must be in YYYY-MM-DD format")
		}
		parsedEnd, err := time.Parse(dateLayout, endDateStr)
		if err != nil {
			return nil, nil, fmt.Errorf("end_date must be in YYYY-MM-DD format")
		}
		if parsedStart.After(parsedEnd) {
			return nil, nil, fmt.Errorf("start_date must not be after end_date")
		}
		return &parsedStart, &parsedEnd, nil
	default:
		return nil, nil, fmt.Errorf("period must be one of 'week', 'month', '3months', 'custom'")
	}
}

func formatValidationError(err error) string {
	validationErrors, ok := err.(validator.ValidationErrors)
	if !ok || len(validationErrors) == 0 {
		return "invalid request"
	}

	field := validationErrors[0]
	switch field.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", field.Field())
	case "gt":
		return fmt.Sprintf("%s must be greater than %s", field.Field(), field.Param())
	case "oneof":
		return fmt.Sprintf("%s must be one of 'Groceries','Leisure','Electronics','Utilities','Clothing','Health','Others'", field.Field())
	default:
		return fmt.Sprintf("%s is invalid", field.Field())
	}
}
