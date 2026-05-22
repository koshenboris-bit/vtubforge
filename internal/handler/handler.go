package handler

import (
	"backend-service/internal/config"
	"backend-service/internal/service"
)

type Handler struct {
	Auth    *AuthHandler
	News    *NewsHandler
	Lessons *LessonHandler
	Users   *UserHandler
}

func NewHandler(cfg config.Config, svcs *service.Services) *Handler {
	return &Handler{
		Auth:    NewAuthHandler(cfg, svcs.Auth),
		News:    NewNewsHandler(svcs.News),
		Lessons: NewLessonHandler(svcs.Lessons),
		Users:   NewUserHandler(svcs.Users),
	}
}
