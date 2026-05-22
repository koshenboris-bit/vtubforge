package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"backend-service/internal/domain"
	"backend-service/internal/service"

	"github.com/go-chi/chi/v5"
)

type NewsHandler struct {
	svc *service.NewsService
}

func NewNewsHandler(svc *service.NewsService) *NewsHandler {
	return &NewsHandler{svc: svc}
}

// List godoc
//
//	@Summary		Get news
//	@Description	Get all news
//	@Tags			news
//	@Security		BearerAuth
//	@Produce		json
//	@Success		200	{array}		domain.News
//	@Failure		401	{object}	map[string]interface{}
//	@Router			/news [get]
func (h *NewsHandler) List(w http.ResponseWriter, r *http.Request) {
	items, err := h.svc.List(r.Context())
	if err != nil {
		writeDomainError(w, err)
		return
	}

	jsonResponse(w, http.StatusOK, items)
}

// Create godoc
//
//	@Summary		Create news
//	@Description	Create new news item
//	@Tags			news
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			request	body		domain.NewsRequest	true	"News payload"
//	@Success		201		{object}	domain.News
//	@Failure		400		{object}	map[string]interface{}
//	@Failure		401		{object}	map[string]interface{}
//	@Failure		403		{object}	map[string]interface{}
//	@Router			/news [post]
func (h *NewsHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req domain.NewsRequest

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

// Delete godoc
//
//	@Summary		Delete news
//	@Description	Delete news by id
//	@Tags			news
//	@Security		BearerAuth
//	@Produce		json
//	@Param			id	path	int	true	"News ID"
//	@Success		204
//	@Failure		400	{object}	map[string]interface{}
//	@Failure		401	{object}	map[string]interface{}
//	@Failure		403	{object}	map[string]interface{}
//	@Failure		404	{object}	map[string]interface{}
//	@Router			/news/{id} [delete]
func (h *NewsHandler) Delete(w http.ResponseWriter, r *http.Request) {
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