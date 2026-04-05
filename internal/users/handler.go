package users

import (
	"encoding/json"
	"net/http"

	"github.com/aadithyaa9/finance-dashboard/internal/response"
	"github.com/go-chi/chi/v5"
)

type Handler struct {
	store *Store
}

func NewHandler(store *Store) *Handler {
	return &Handler{store: store}
}

// List godoc
// GET /api/users  (admin only)
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	list, err := h.store.List()
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "failed to fetch users")
		return
	}
	response.JSON(w, http.StatusOK, list)
}

// UpdateRole godoc
// PATCH /api/users/{id}/role  (admin only)
func (h *Handler) UpdateRole(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var body struct {
		Role string `json:"role"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	switch Role(body.Role) {
	case RoleViewer, RoleAnalyst, RoleAdmin:
	default:
		response.Error(w, http.StatusBadRequest, "role must be one of: viewer, analyst, admin")
		return
	}

	if err := h.store.UpdateRole(id, Role(body.Role)); err != nil {
		response.Error(w, http.StatusInternalServerError, "failed to update role")
		return
	}

	response.Message(w, http.StatusOK, "role updated")
}

// UpdateStatus godoc
// PATCH /api/users/{id}/status  (admin only)
func (h *Handler) UpdateStatus(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var body struct {
		IsActive bool `json:"is_active"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := h.store.UpdateStatus(id, body.IsActive); err != nil {
		response.Error(w, http.StatusInternalServerError, "failed to update status")
		return
	}

	response.Message(w, http.StatusOK, "status updated")
}
