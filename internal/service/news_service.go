package service

import (
	"context"

	"backend-service/internal/domain"
	"backend-service/internal/repository"
)

type NewsService struct {
	repo *repository.NewsRepository
}

func NewNewsService(repos *repository.Repositories) *NewsService {
	return &NewsService{repo: repos.News}
}

func (s *NewsService) Create(ctx context.Context, req domain.NewsRequest) (domain.News, error) {
	if req.Title == "" || req.Content == "" {
		return domain.News{}, domain.ErrValidation
	}
	return s.repo.Create(ctx, domain.News{Title: req.Title, Content: req.Content})
}

func (s *NewsService) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}

func (s *NewsService) List(ctx context.Context) ([]domain.News, error) {
	return s.repo.List(ctx)
}
