package records

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/aadithyaa9/finance-dashboard/internal/middleware"
	"github.com/aadithyaa9/finance-dashboard/internal/response"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

type Handler struct {
	store *Store
}

func NewHandler(store *Store) *Handler {
	return &Handler{store: store}
}

// List godoc
// GET /api/records?type=income&category=food&from=2024-01-01&to=2024-12-31&page=1&limit=20
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	page, _ := strconv.Atoi(q.Get("page"))
	limit, _ := strconv.Atoi(q.Get("limit"))

	filter := Filter{
		Type:     q.Get("type"),
		Category: q.Get("category"),
		From:     q.Get("from"),
		To:       q.Get("to"),
		Page:     page,
		Limit:    limit,
	}

	list, total, err := h.store.List(filter)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "failed to fetch records")
		return
	}

	response.JSON(w, http.StatusOK, map[string]interface{}{
		"records": list,
		"total":   total,
		"page":    filter.Page,
		"limit":   filter.Limit,
	})
}

// Create godoc
// POST /api/records
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	claims := middleware.ClaimsFromContext(r.Context())

	var input CreateInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if err := validate.Struct(input); err != nil {
		response.Error(w, http.StatusUnprocessableEntity, validationMessage(err))
		return
	}

	rec, err := h.store.Create(claims.UserID, input)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	response.JSON(w, http.StatusCreated, rec)
}

// GetByID godoc
// GET /api/records/{id}
func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	rec, err := h.store.GetByID(id)
	if err != nil {
		response.Error(w, http.StatusNotFound, "record not found")
		return
	}
	response.JSON(w, http.StatusOK, rec)
}

// Update godoc
// PUT /api/records/{id}
func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var input UpdateInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if err := validate.Struct(input); err != nil {
		response.Error(w, http.StatusUnprocessableEntity, validationMessage(err))
		return
	}

	rec, err := h.store.Update(id, input)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	response.JSON(w, http.StatusOK, rec)
}

// Delete godoc
// DELETE /api/records/{id}
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := h.store.SoftDelete(id); err != nil {
		response.Error(w, http.StatusNotFound, err.Error())
		return
	}
	response.Message(w, http.StatusOK, "record deleted")
}

func validationMessage(err error) string {
	var ve validator.ValidationErrors
	if errors.As(err, &ve) {
		out := ""
		for i, e := range ve {
			if i > 0 {
				out += "; "
			}
			out += e.Field() + ": " + e.Tag()
		}
		return out
	}
	return err.Error()
}
