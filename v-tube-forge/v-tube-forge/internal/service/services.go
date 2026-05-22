package service

import (
	"backend-service/internal/config"
	"backend-service/internal/repository"
)

type Services struct {
	Auth   *AuthService
	News   *NewsService
	Lessons *LessonService
	Users  *UserService
}

func NewServices(cfg config.Config, repos *repository.Repositories) *Services {
	return &Services{
		Auth:    NewAuthService(cfg, repos),
		News:    NewNewsService(repos),
		Lessons: NewLessonService(repos),
		Users:   NewUserService(repos),
	}
}
