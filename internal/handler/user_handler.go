package handler

import (
	"net/http"
	"strconv"

	"backend-service/internal/middleware"
	"backend-service/internal/service"

	"github.com/go-chi/chi/v5"
)

type UserHandler struct {
	svc *service.UserService
}

func NewUserHandler(svc *service.UserService) *UserHandler {
	return &UserHandler{svc: svc}
}

// List users godoc
// @Summary Get users
// @Tags users
// @Security BearerAuth
// @Produce json
// @Success 200 {array} domain.User
// @Router /users [get]
func (h *UserHandler) List(w http.ResponseWriter, r *http.Request) {
	items, err := h.svc.List(r.Context())
	if err != nil {
		writeDomainError(w, err)
		return
	}
	jsonResponse(w, http.StatusOK, items)
}

// Get user godoc
// @Summary Get user by id
// @Tags users
// @Security BearerAuth
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} domain.User
// @Router /users/{id} [get]
func (h *UserHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		jsonError(w, http.StatusBadRequest, "invalid id")
		return
	}
	item, err := h.svc.GetByID(r.Context(), id)
	if err != nil {
		writeDomainError(w, err)
		return
	}
	jsonResponse(w, http.StatusOK, item)
}

// Me godoc
// @Summary Get current user
// @Tags users
// @Security BearerAuth
// @Produce json
// @Success 200 {object} domain.User
// @Router /users/me [get]
func (h *UserHandler) Me(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserID(r.Context())
	if !ok {
		jsonError(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	item, err := h.svc.Me(r.Context(), userID)
	if err != nil {
		writeDomainError(w, err)
		return
	}
	jsonResponse(w, http.StatusOK, item)
}
