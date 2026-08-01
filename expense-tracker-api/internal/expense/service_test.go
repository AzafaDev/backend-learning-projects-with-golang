package expense

import (
	"context"
	"expense-tracker-api/internal/httpserver"
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCalculatePeriodRange(t *testing.T) {
	now := time.Date(2026, time.August, 15, 13, 45, 0, 0, time.UTC)
	today := time.Date(2026, time.August, 15, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name      string
		period    string
		startStr  string
		endStr    string
		wantStart *time.Time
		wantEnd   *time.Time
		wantErr   bool
	}{
		{
			name:      "no period means no filter",
			period:    "",
			wantStart: nil,
			wantEnd:   nil,
		},
		{
			name:      "week",
			period:    "week",
			wantStart: ptr(today.AddDate(0, 0, -7)),
			wantEnd:   ptr(today),
		},
		{
			name:      "month",
			period:    "month",
			wantStart: ptr(today.AddDate(0, -1, 0)),
			wantEnd:   ptr(today),
		},
		{
			name:      "3months",
			period:    "3months",
			wantStart: ptr(today.AddDate(0, -3, 0)),
			wantEnd:   ptr(today),
		},
		{
			name:      "custom valid range",
			period:    "custom",
			startStr:  "2026-01-01",
			endStr:    "2026-01-31",
			wantStart: ptr(time.Date(2026, time.January, 1, 0, 0, 0, 0, time.UTC)),
			wantEnd:   ptr(time.Date(2026, time.January, 31, 0, 0, 0, 0, time.UTC)),
		},
		{
			name:     "custom missing start_date",
			period:   "custom",
			endStr:   "2026-01-31",
			wantErr:  true,
		},
		{
			name:     "custom missing end_date",
			period:   "custom",
			startStr: "2026-01-01",
			wantErr:  true,
		},
		{
			name:     "custom start after end",
			period:   "custom",
			startStr: "2026-01-31",
			endStr:   "2026-01-01",
			wantErr:  true,
		},
		{
			name:     "custom bad date format",
			period:   "custom",
			startStr: "01/01/2026",
			endStr:   "2026-01-31",
			wantErr:  true,
		},
		{
			name:    "unknown period",
			period:  "year",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			start, end, err := calculatePeriodRange(tt.period, tt.startStr, tt.endStr, now)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.wantStart, start)
			assert.Equal(t, tt.wantEnd, end)
		})
	}
}

func ptr(t time.Time) *time.Time {
	return &t
}

type mockExpenseRepo struct {
	mock.Mock
}

func (m *mockExpenseRepo) GetExpenseByID(ctx context.Context, id, userID uuid.UUID) (*Expense, error) {
	args := m.Called(ctx, id, userID)
	expense, _ := args.Get(0).(*Expense)
	return expense, args.Error(1)
}

func (m *mockExpenseRepo) GetExpenses(ctx context.Context, userID uuid.UUID, filter ExpenseFilter, pagination Pagination) ([]Expense, int, error) {
	args := m.Called(ctx, userID, filter, pagination)
	expenses, _ := args.Get(0).([]Expense)
	return expenses, args.Int(1), args.Error(2)
}

func (m *mockExpenseRepo) CreateExpense(ctx context.Context, userID uuid.UUID, req CreateExpenseRequest) (*Expense, error) {
	args := m.Called(ctx, userID, req)
	expense, _ := args.Get(0).(*Expense)
	return expense, args.Error(1)
}

func (m *mockExpenseRepo) UpdateExpense(ctx context.Context, id, userID uuid.UUID, req UpdateExpenseRequest) (*Expense, error) {
	args := m.Called(ctx, id, userID, req)
	expense, _ := args.Get(0).(*Expense)
	return expense, args.Error(1)
}

func (m *mockExpenseRepo) DeleteExpense(ctx context.Context, id, userID uuid.UUID) error {
	args := m.Called(ctx, id, userID)
	return args.Error(0)
}

func validCreateRequest() CreateExpenseRequest {
	return CreateExpenseRequest{
		Title:    "Groceries run",
		Amount:   50000,
		Category: Groceries,
		Date:     time.Now(),
	}
}

func TestExpenseService_CreateExpense_ValidationError(t *testing.T) {
	repo := new(mockExpenseRepo)
	svc := NewExpenseService(repo)
	userID := uuid.New()

	req := validCreateRequest()
	req.Category = "NotACategory"

	_, err := svc.CreateExpense(context.Background(), userID, req)

	assert.Error(t, err)
	var clientErr *httpserver.ClientError
	assert.ErrorAs(t, err, &clientErr)
	assert.Equal(t, http.StatusBadRequest, clientErr.Status)
	repo.AssertNotCalled(t, "CreateExpense")
}

func TestExpenseService_CreateExpense_Success(t *testing.T) {
	repo := new(mockExpenseRepo)
	svc := NewExpenseService(repo)
	userID := uuid.New()
	req := validCreateRequest()

	want := &Expense{ID: uuid.New(), UserID: userID, Title: req.Title, Amount: req.Amount, Category: req.Category, Date: req.Date}
	repo.On("CreateExpense", mock.Anything, userID, req).Return(want, nil)

	got, err := svc.CreateExpense(context.Background(), userID, req)

	assert.NoError(t, err)
	assert.Equal(t, want, got)
	repo.AssertExpectations(t)
}

func TestExpenseService_GetExpenseByID_NotFoundPassesThrough(t *testing.T) {
	repo := new(mockExpenseRepo)
	svc := NewExpenseService(repo)
	userID := uuid.New()
	id := uuid.New()

	notFound := httpserver.NewClientError(http.StatusNotFound, "expense not found")
	repo.On("GetExpenseByID", mock.Anything, id, userID).Return(nil, notFound)

	_, err := svc.GetExpenseByID(context.Background(), userID, id)

	assert.Error(t, err)
	var clientErr *httpserver.ClientError
	assert.ErrorAs(t, err, &clientErr)
	assert.Equal(t, http.StatusNotFound, clientErr.Status)
	repo.AssertExpectations(t)
}
