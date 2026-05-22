package handler

import (
	"encoding/json"
	"net/http"

	"backend-service/internal/config"
	"backend-service/internal/domain"
	"backend-service/internal/service"
)

type AuthHandler struct {
	cfg  config.Config
	svc  *service.AuthService
}

func NewAuthHandler(cfg config.Config, svc *service.AuthService) *AuthHandler {
	return &AuthHandler{cfg: cfg, svc: svc}
}

// Register godoc
// @Summary Register user
// @Tags auth
// @Accept json
// @Produce json
// @Param request body domain.RegisterRequest true "Register request"
// @Success 201 {object} domain.AuthResponse
// @Failure 400 {object} map[string]interface{}
// @Router /register [post]
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req domain.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, http.StatusBadRequest, "invalid json")
		return
	}
	resp, err := h.svc.Register(r.Context(), req)
	if err != nil {
		writeDomainError(w, err)
		return
	}
	jsonResponse(w, http.StatusCreated, resp)
}

// Login godoc
// @Summary Login user
// @Tags auth
// @Accept json
// @Produce json
// @Param request body domain.LoginRequest true "Login request"
// @Success 200 {object} domain.AuthResponse
// @Failure 401 {object} map[string]interface{}
// @Router /login [post]
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req domain.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, http.StatusBadRequest, "invalid json")
		return
	}
	resp, err := h.svc.Login(r.Context(), req)
	if err != nil {
		writeDomainError(w, err)
		return
	}
	jsonResponse(w, http.StatusOK, resp)
}

// Refresh godoc
// @Summary Refresh access token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body domain.RefreshRequest true "Refresh request"
// @Success 200 {object} domain.TokenPair
// @Router /refresh [post]
func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	var req domain.RefreshRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, http.StatusBadRequest, "invalid json")
		return
	}
	pair, err := h.svc.Refresh(r.Context(), req.RefreshToken)
	if err != nil {
		writeDomainError(w, err)
		return
	}
	jsonResponse(w, http.StatusOK, pair)
}
