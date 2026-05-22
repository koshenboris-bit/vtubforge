package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"backend-service/internal/domain"
	"backend-service/internal/middleware"
	"backend-service/internal/service"

	"github.com/go-chi/chi/v5"
)

type LessonHandler struct {
	svc *service.LessonService
}

func NewLessonHandler(svc *service.LessonService) *LessonHandler {
	return &LessonHandler{svc: svc}
}

// List lessons godoc
// @Summary Get lessons
// @Tags lessons
// @Security BearerAuth
// @Produce json
// @Success 200 {array} domain.Lesson
// @Router /lessons [get]
func (h *LessonHandler) List(w http.ResponseWriter, r *http.Request) {
	callerID, _ := middleware.UserID(r.Context())

	items, err := h.svc.List(r.Context(), callerID)
	if err != nil {
		writeDomainError(w, err)
		return
	}
	jsonResponse(w, http.StatusOK, items)
}

// Create lesson godoc
// @Summary Create lesson
// @Tags lessons
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body domain.LessonRequest true "Lesson"
// @Success 201 {object} domain.Lesson
// @Router /lessons [post]
func (h *LessonHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req domain.LessonRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, http.StatusBadRequest, "invalid json")
		return
	}
	item, err := h.svc.Create(r.Context(), req)
	if err != nil {
		writeDomainError(w, err)
		return
	}
	jsonResponse(w, http.StatusCreated, item)
}

// Update lesson godoc
// @Summary Update lesson
// @Tags lessons
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Lesson ID"
// @Param request body domain.LessonRequest true "Lesson"
// @Success 200 {object} domain.Lesson
// @Router /lessons/{id} [put]
func (h *LessonHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		jsonError(w, http.StatusBadRequest, "invalid id")
		return
	}
	var req domain.LessonRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, http.StatusBadRequest, "invalid json")
		return
	}
	item, err := h.svc.Update(r.Context(), id, req)
	if err != nil {
		writeDomainError(w, err)
		return
	}
	jsonResponse(w, http.StatusOK, item)
}

// Delete lesson godoc
// @Summary Delete lesson
// @Tags lessons
// @Security BearerAuth
// @Param id path int true "Lesson ID"
// @Success 204
// @Router /lessons/{id} [delete]
func (h *LessonHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		jsonError(w, http.StatusBadRequest, "invalid id")
		return
	}
	if err := h.svc.Delete(r.Context(), id); err != nil {
		writeDomainError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// Pass lesson godoc
// @Summary Pass lesson
// @Tags lessons
// @Security BearerAuth
// @Param lessonId path int true "Lesson ID"
// @Param userId path int true "User ID"
// @Success 204
// @Router /lessons/{lessonId}/pass/{userId} [post]
func (h *LessonHandler) Pass(w http.ResponseWriter, r *http.Request) {
	lessonID, err := strconv.ParseInt(chi.URLParam(r, "lessonId"), 10, 64)
	if err != nil {
		jsonError(w, http.StatusBadRequest, "invalid lesson id")
		return
	}
	userID, err := strconv.ParseInt(chi.URLParam(r, "userId"), 10, 64)
	if err != nil {
		jsonError(w, http.StatusBadRequest, "invalid user id")
		return
	}
	callerID, _ := middleware.UserID(r.Context())
	callerRole, _ := middleware.UserRole(r.Context())
	if err := h.svc.Pass(r.Context(), lessonID, userID, callerID, callerRole); err != nil {
		writeDomainError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
