package auth

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/aadithyaa9/finance-dashboard/internal/response"
	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

// Register godoc
// POST /api/auth/register
func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var input RegisterInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if err := validate.Struct(input); err != nil {
		response.Error(w, http.StatusUnprocessableEntity, validationMessage(err))
		return
	}

	res, err := h.svc.Register(input)
	if err != nil {
		switch {
		case errors.Is(err, ErrEmailTaken):
			response.Error(w, http.StatusConflict, err.Error())
		default:
			response.Error(w, http.StatusInternalServerError, "registration failed")
		}
		return
	}

	response.JSON(w, http.StatusCreated, res)
}

// Login godoc
// POST /api/auth/login
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var input LoginInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if err := validate.Struct(input); err != nil {
		response.Error(w, http.StatusUnprocessableEntity, validationMessage(err))
		return
	}

	res, err := h.svc.Login(input)
	if err != nil {
		switch {
		case errors.Is(err, ErrInvalidCredentials):
			response.Error(w, http.StatusUnauthorized, err.Error())
		case errors.Is(err, ErrInactiveUser):
			response.Error(w, http.StatusForbidden, err.Error())
		default:
			response.Error(w, http.StatusInternalServerError, "login failed")
		}
		return
	}

	response.JSON(w, http.StatusOK, res)
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
