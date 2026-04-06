package dashboard

import (
	"net/http"

	"github.com/aadithyaa9/finance-dashboard/internal/middleware"
	"github.com/aadithyaa9/finance-dashboard/internal/response"
)

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) Summary(w http.ResponseWriter, r *http.Request) {
	claims := middleware.ClaimsFromContext(r.Context())
	data, err := h.svc.GetSummary(claims.UserID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "failed to fetch summary")
		return
	}
	response.JSON(w, http.StatusOK, data)
}

func (h *Handler) ByCategory(w http.ResponseWriter, r *http.Request) {
	claims := middleware.ClaimsFromContext(r.Context())
	data, err := h.svc.GetByCategory(claims.UserID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "failed to fetch category totals")
		return
	}
	response.JSON(w, http.StatusOK, data)
}

func (h *Handler) Trends(w http.ResponseWriter, r *http.Request) {
	claims := middleware.ClaimsFromContext(r.Context())
	data, err := h.svc.GetMonthlyTrends(claims.UserID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "failed to fetch trends")
		return
	}
	response.JSON(w, http.StatusOK, data)
}

func (h *Handler) Recent(w http.ResponseWriter, r *http.Request) {
	claims := middleware.ClaimsFromContext(r.Context())
	data, err := h.svc.GetRecent(claims.UserID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "failed to fetch recent records")
		return
	}
	response.JSON(w, http.StatusOK, data)
}