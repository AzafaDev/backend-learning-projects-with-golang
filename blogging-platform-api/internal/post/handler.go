package post

import (
	"blogging-platform-api/internal/httpserver"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/google/uuid"
)

type Handler struct {
	srv *Service
}

func NewHandler(srv *Service) *Handler {
	return &Handler{
		srv: srv,
	}
}

func (h *Handler) CreatePost(w http.ResponseWriter, r *http.Request) {
	var req CreatePostRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpserver.WriteError(w, http.StatusBadRequest, "invalid request body", err)
		return
	}

	createdPost, err := h.srv.CreatePost(r.Context(), req)
	if err != nil {
		if errors.Is(err, ErrBadRequest) {
			httpserver.WriteError(w, http.StatusBadRequest, err.Error(), err)
			return
		}
		httpserver.WriteError(w, http.StatusInternalServerError, "something went wrong", err)
		return
	}

	httpserver.WriteJSON(w, http.StatusCreated, createdPost)
}

func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		httpserver.WriteError(w, http.StatusBadRequest, "id is invalid", err)
		return
	}

	existingPost, err := h.srv.GetPostByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, ErrPostNotFound) {
			httpserver.WriteError(w, http.StatusNotFound, err.Error(), err)
			return
		}
		httpserver.WriteError(w, http.StatusInternalServerError, "something went wrong", err)
		return
	}

	httpserver.WriteJSON(w, http.StatusOK, existingPost)
}

func (h *Handler) UpdatePost(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		httpserver.WriteError(w, http.StatusBadRequest, "id is invalid", err)
		return
	}

	var req UpdatePostRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpserver.WriteError(w, http.StatusBadRequest, "invalid request body", err)
		return
	}

	updatedPost, err := h.srv.UpdatePost(r.Context(), req, id)
	if err != nil {
		if errors.Is(err, ErrPostNotFound) {
			httpserver.WriteError(w, http.StatusNotFound, err.Error(), err)
			return
		}
		if errors.Is(err, ErrBadRequest) {
			httpserver.WriteError(w, http.StatusBadRequest, err.Error(), err)
			return
		}
		httpserver.WriteError(w, http.StatusInternalServerError, "something went wrong", err)
		return
	}

	httpserver.WriteJSON(w, http.StatusOK, updatedPost)
}

func (h *Handler) DeletePost(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		httpserver.WriteError(w, http.StatusBadRequest, "id is invalid", err)
		return
	}

	_, err = h.srv.DeletePostByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, ErrPostNotFound) {
			httpserver.WriteError(w, http.StatusNotFound, err.Error(), err)
			return
		}
		httpserver.WriteError(w, http.StatusInternalServerError, "something went wrong", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) SearchPost(w http.ResponseWriter, r *http.Request) {
	term := r.URL.Query().Get("term")

	listPosts, err := h.srv.SearchPosts(r.Context(), term)
	if err != nil {
		httpserver.WriteError(w, http.StatusInternalServerError, "something went wrong", err)
		return
	}

	httpserver.WriteJSON(w, http.StatusOK, listPosts)
}
