package handler

import (
	"encoding/json"
	"net/http"

	"backend-service/internal/domain"
)

func jsonResponse(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func jsonError(w http.ResponseWriter, status int, msg string) {
	jsonResponse(w, status, map[string]any{"error": msg})
}

func writeDomainError(w http.ResponseWriter, err error) {
	switch err {
	case domain.ErrValidation:
		jsonError(w, http.StatusBadRequest, err.Error())
	case domain.ErrUnauthorized, domain.ErrInvalidToken, domain.ErrPasswordMismatch:
		jsonError(w, http.StatusUnauthorized, err.Error())
	case domain.ErrForbidden:
		jsonError(w, http.StatusForbidden, err.Error())
	case domain.ErrNotFound:
		jsonError(w, http.StatusNotFound, err.Error())
	case domain.ErrConflict:
		jsonError(w, http.StatusConflict, err.Error())
	default:
		jsonError(w, http.StatusInternalServerError, "internal server error")
	}
}
