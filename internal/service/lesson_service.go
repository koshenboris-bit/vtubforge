package service

import (
	"context"

	"backend-service/internal/domain"
	"backend-service/internal/repository"
)

type LessonService struct {
	repo *repository.LessonRepository
	users *repository.UserRepository
}

func NewLessonService(repos *repository.Repositories) *LessonService {
	return &LessonService{repo: repos.Lessons, users: repos.Users}
}

func (s *LessonService) Create(ctx context.Context, req domain.LessonRequest) (domain.Lesson, error) {
	if req.Title == "" || req.Description == "" || req.VideoLink == "" || req.LessonType == "" {
		return domain.Lesson{}, domain.ErrValidation
	}
	if req.LessonType != domain.LessonIntro && req.LessonType != domain.LessonFull {
		return domain.Lesson{}, domain.ErrValidation
	}
	return s.repo.Create(ctx, domain.Lesson{
		Title: req.Title, Description: req.Description, LessonType: req.LessonType, VideoLink: req.VideoLink,
	})
}

func (s *LessonService) Update(ctx context.Context, id int64, req domain.LessonRequest) (domain.Lesson, error) {
	if req.Title == "" || req.Description == "" || req.VideoLink == "" || req.LessonType == "" {
		return domain.Lesson{}, domain.ErrValidation
	}
	return s.repo.Update(ctx, id, domain.Lesson{
		Title: req.Title, Description: req.Description, LessonType: req.LessonType, VideoLink: req.VideoLink,
	})
}

func (s *LessonService) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}

func (s *LessonService) List(ctx context.Context, userID int64) ([]domain.Lesson, error) {
	return s.repo.List(ctx, userID)
}

func (s *LessonService) Pass(ctx context.Context, lessonID, userID int64, callerID int64, callerRole string) error {
	if callerRole != "admin" && callerID != userID {
		return domain.ErrForbidden
	}
	_, err := s.repo.GetByID(ctx, lessonID)
	if err != nil {
		return err
	}
	_, err = s.users.GetByID(ctx, userID)
	if err != nil {
		return err
	}
	return s.repo.MarkPassed(ctx, userID, lessonID)
}
