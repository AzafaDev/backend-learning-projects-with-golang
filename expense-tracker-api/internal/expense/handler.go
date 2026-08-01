package expense

import (
	"encoding/json"
	"errors"
	"expense-tracker-api/internal/httpserver"
	"expense-tracker-api/internal/middleware"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type ExpenseHandler struct {
	srv *ExpenseService
}

func NewExpenseHandler(srv *ExpenseService) *ExpenseHandler {
	return &ExpenseHandler{
		srv: srv,
	}
}

func (h *ExpenseHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID, ok := userIDFromRequest(w, r)
	if !ok {
		return
	}

	var req CreateExpenseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Error().Err(err).Msg("create: expense handler")
		httpserver.WriteJSONError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	createdExpense, err := h.srv.CreateExpense(r.Context(), userID, req)
	if err != nil {
		httpserver.RespondError(w, "create: expense handler", err)
		return
	}

	httpserver.WriteJSON(w, createdExpense, http.StatusCreated)
}

func (h *ExpenseHandler) List(w http.ResponseWriter, r *http.Request) {
	userID, ok := userIDFromRequest(w, r)
	if !ok {
		return
	}

	query := r.URL.Query()
	q := ListExpensesQuery{
		Period:    query.Get("period"),
		StartDate: query.Get("start_date"),
		EndDate:   query.Get("end_date"),
		Category:  query.Get("category"),
	}
	if page, err := strconv.Atoi(query.Get("page")); err == nil {
		q.Page = page
	}
	if limit, err := strconv.Atoi(query.Get("limit")); err == nil {
		q.Limit = limit
	}

	result, err := h.srv.GetExpenses(r.Context(), userID, q)
	if err != nil {
		httpserver.RespondError(w, "list: expense handler", err)
		return
	}

	httpserver.WriteJSON(w, result, http.StatusOK)
}

func (h *ExpenseHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	userID, ok := userIDFromRequest(w, r)
	if !ok {
		return
	}

	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		httpserver.WriteJSONError(w, "invalid expense id", http.StatusBadRequest)
		return
	}

	expense, err := h.srv.GetExpenseByID(r.Context(), userID, id)
	if err != nil {
		httpserver.RespondError(w, "get by id: expense handler", err)
		return
	}

	httpserver.WriteJSON(w, expense, http.StatusOK)
}

func (h *ExpenseHandler) Update(w http.ResponseWriter, r *http.Request) {
	userID, ok := userIDFromRequest(w, r)
	if !ok {
		return
	}

	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		httpserver.WriteJSONError(w, "invalid expense id", http.StatusBadRequest)
		return
	}

	var req UpdateExpenseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Error().Err(err).Msg("update: expense handler")
		httpserver.WriteJSONError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	updatedExpense, err := h.srv.UpdateExpense(r.Context(), userID, id, req)
	if err != nil {
		httpserver.RespondError(w, "update: expense handler", err)
		return
	}

	httpserver.WriteJSON(w, updatedExpense, http.StatusOK)
}

func (h *ExpenseHandler) Delete(w http.ResponseWriter, r *http.Request) {
	userID, ok := userIDFromRequest(w, r)
	if !ok {
		return
	}

	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		httpserver.WriteJSONError(w, "invalid expense id", http.StatusBadRequest)
		return
	}

	if err := h.srv.DeleteExpense(r.Context(), userID, id); err != nil {
		httpserver.RespondError(w, "delete: expense handler", err)
		return
	}

	httpserver.WriteJSON(w, map[string]string{"message": "expense deleted"}, http.StatusOK)
}

func userIDFromRequest(w http.ResponseWriter, r *http.Request) (uuid.UUID, bool) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		httpserver.RespondError(w, "expense handler", errors.New("missing user id in request context"))
		return uuid.UUID{}, false
	}
	return userID, true
}
