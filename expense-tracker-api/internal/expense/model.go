package expense

import (
	"time"

	"github.com/google/uuid"
)

type Expense struct {
	ID          uuid.UUID       `json:"id"`
	UserID      uuid.UUID       `json:"user_id"`
	Title       string          `json:"title"`
	Amount      int64           `json:"amount"`
	Category    ExpenseCategory `json:"category"`
	Description *string         `json:"description,omitempty"`
	Date        time.Time       `json:"date"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
}

type ExpenseCategory string

const (
	Groceries   ExpenseCategory = "Groceries"
	Leisure     ExpenseCategory = "Leisure"
	Electronics ExpenseCategory = "Electronics"
	Utilities   ExpenseCategory = "Utilities"
	Clothing    ExpenseCategory = "Clothing"
	Health      ExpenseCategory = "Health"
	Others      ExpenseCategory = "Others"
)

type CreateExpenseRequest struct {
	Title       string          `json:"title" validate:"required"`
	Amount      int64           `json:"amount" validate:"required,gt=0"`
	Category    ExpenseCategory `json:"category" validate:"required,oneof=Groceries Leisure Electronics Utilities Clothing Health Others"`
	Description *string         `json:"description,omitempty"`
	Date        time.Time       `json:"date" validate:"required"`
}

type UpdateExpenseRequest struct {
	Title       string          `json:"title" validate:"required"`
	Amount      int64           `json:"amount" validate:"required,gt=0"`
	Category    ExpenseCategory `json:"category" validate:"required,oneof=Groceries Leisure Electronics Utilities Clothing Health Others"`
	Description *string         `json:"description,omitempty"`
	Date        time.Time       `json:"date" validate:"required"`
}

type ExpenseFilter struct {
	StartDate *time.Time
	EndDate   *time.Time
	Category  *ExpenseCategory
}

type Pagination struct {
	Page  int
	Limit int
}

type ExpenseListResponse struct {
	Data  []Expense `json:"data"`
	Page  int       `json:"page"`
	Limit int       `json:"limit"`
	Total int       `json:"total"`
}

