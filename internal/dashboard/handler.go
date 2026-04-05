package dashboard

import (
	"net/http"

	"github.com/aadithyaa9/finance-dashboard/internal/response"
)

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

// Summary godoc
// GET /api/dashboard/summary
func (h *Handler) Summary(w http.ResponseWriter, r *http.Request) {
	data, err := h.svc.GetSummary()
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "failed to fetch summary")
		return
	}
	response.JSON(w, http.StatusOK, data)
}

// ByCategory godoc
// GET /api/dashboard/by-category
func (h *Handler) ByCategory(w http.ResponseWriter, r *http.Request) {
	data, err := h.svc.GetByCategory()
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "failed to fetch category totals")
		return
	}
	response.JSON(w, http.StatusOK, data)
}

// Trends godoc
// GET /api/dashboard/trends
func (h *Handler) Trends(w http.ResponseWriter, r *http.Request) {
	data, err := h.svc.GetMonthlyTrends()
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "failed to fetch trends")
		return
	}
	response.JSON(w, http.StatusOK, data)
}

// Recent godoc
// GET /api/dashboard/recent
func (h *Handler) Recent(w http.ResponseWriter, r *http.Request) {
	data, err := h.svc.GetRecent()
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "failed to fetch recent records")
		return
	}
	response.JSON(w, http.StatusOK, data)
}
