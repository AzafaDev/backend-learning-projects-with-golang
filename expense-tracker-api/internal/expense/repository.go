package expense

import (
	"context"
	"errors"
	"expense-tracker-api/internal/httpserver"
	"fmt"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ExpenseRepository struct {
	db *pgxpool.Pool
}

func NewExpenseRepository(db *pgxpool.Pool) *ExpenseRepository {
	return &ExpenseRepository{
		db: db,
	}
}

func (e *ExpenseRepository) GetExpenseByID(ctx context.Context, id, userID uuid.UUID) (*Expense, error) {
	var expense Expense
	query := `
	SELECT id, user_id, title, amount, category, description, date, created_at, updated_at
	FROM expenses
	WHERE id=$1 AND user_id=$2
	`

	if err := e.db.QueryRow(ctx, query, id, userID).Scan(
		&expense.ID,
		&expense.UserID,
		&expense.Title,
		&expense.Amount,
		&expense.Category,
		&expense.Description,
		&expense.Date,
		&expense.CreatedAt,
		&expense.UpdatedAt,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, httpserver.NewClientError(http.StatusNotFound, "expense not found")
		}
		return nil, err
	}

	return &expense, nil
}

// GetExpenses builds the WHERE clause dynamically from filter, and returns
// the page of results alongside the total row count (via COUNT(*) OVER()).
func (e *ExpenseRepository) GetExpenses(ctx context.Context, userID uuid.UUID, filter ExpenseFilter, pagination Pagination) ([]Expense, int, error) {
	conditions := []string{"user_id = $1"}
	args := []any{userID}

	if filter.StartDate != nil {
		args = append(args, *filter.StartDate)
		conditions = append(conditions, fmt.Sprintf("date >= $%d", len(args)))
	}
	if filter.EndDate != nil {
		args = append(args, *filter.EndDate)
		conditions = append(conditions, fmt.Sprintf("date <= $%d", len(args)))
	}
	if filter.Category != nil {
		args = append(args, *filter.Category)
		conditions = append(conditions, fmt.Sprintf("category = $%d", len(args)))
	}

	args = append(args, pagination.Limit)
	limitIdx := len(args)
	args = append(args, (pagination.Page-1)*pagination.Limit)
	offsetIdx := len(args)

	query := fmt.Sprintf(`
	SELECT id, user_id, title, amount, category, description, date, created_at, updated_at,
	       COUNT(*) OVER() AS total
	FROM expenses
	WHERE %s
	ORDER BY date DESC
	LIMIT $%d OFFSET $%d
	`, strings.Join(conditions, " AND "), limitIdx, offsetIdx)

	rows, err := e.db.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var expenses []Expense
	total := 0

	for rows.Next() {
		var expense Expense

		if err := rows.Scan(
			&expense.ID,
			&expense.UserID,
			&expense.Title,
			&expense.Amount,
			&expense.Category,
			&expense.Description,
			&expense.Date,
			&expense.CreatedAt,
			&expense.UpdatedAt,
			&total,
		); err != nil {
			return nil, 0, err
		}

		expenses = append(expenses, expense)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return expenses, total, nil
}

func (e *ExpenseRepository) CreateExpense(ctx context.Context, userID uuid.UUID, req CreateExpenseRequest) (*Expense, error) {
	var expense Expense
	query := `
	INSERT INTO expenses
	(user_id, title, amount, category, description, date)
	VALUES ($1, $2, $3, $4, $5, $6)
	RETURNING id, user_id, title, amount, category, description, date, created_at, updated_at
	`

	if err := e.db.QueryRow(ctx, query,
		userID,
		req.Title,
		req.Amount,
		req.Category,
		req.Description,
		req.Date,
	).Scan(
		&expense.ID,
		&expense.UserID,
		&expense.Title,
		&expense.Amount,
		&expense.Category,
		&expense.Description,
		&expense.Date,
		&expense.CreatedAt,
		&expense.UpdatedAt,
	); err != nil {
		return nil, err
	}

	return &expense, nil
}

func (e *ExpenseRepository) UpdateExpense(ctx context.Context, id, userID uuid.UUID, req UpdateExpenseRequest) (*Expense, error) {
	var expense Expense
	query := `
	UPDATE expenses
	SET title=$1, amount=$2, category=$3, description=$4, date=$5, updated_at=now()
	WHERE id=$6 AND user_id=$7
	RETURNING id, user_id, title, amount, category, description, date, created_at, updated_at
	`

	if err := e.db.QueryRow(ctx, query,
		req.Title,
		req.Amount,
		req.Category,
		req.Description,
		req.Date,
		id,
		userID,
	).Scan(
		&expense.ID,
		&expense.UserID,
		&expense.Title,
		&expense.Amount,
		&expense.Category,
		&expense.Description,
		&expense.Date,
		&expense.CreatedAt,
		&expense.UpdatedAt,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, httpserver.NewClientError(http.StatusNotFound, "expense not found")
		}
		return nil, err
	}

	return &expense, nil
}

func (e *ExpenseRepository) DeleteExpense(ctx context.Context, id, userID uuid.UUID) error {
	query := `DELETE FROM expenses WHERE id=$1 AND user_id=$2`

	tag, err := e.db.Exec(ctx, query, id, userID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return httpserver.NewClientError(http.StatusNotFound, "expense not found")
	}

	return nil
}
