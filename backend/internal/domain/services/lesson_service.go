package services

import (
	"context"
	"gitlab.com/w0ikid/study-platform/internal/domain/models"
	"gitlab.com/w0ikid/study-platform/internal/domain/repositories"
)

type LessonServiceInterface interface {
	CreateLesson(ctx context.Context, lesson *models.Lesson) (*models.Lesson, error)
	GetLessonByID(ctx context.Context, id int) (*models.Lesson, error)
	GetAllLessons(ctx context.Context, courseID int) ([]*models.Lesson, error)
	UpdateLesson(ctx context.Context, lesson *models.Lesson) error
	DeleteLesson(ctx context.Context, id int) error
}

type LessonService struct {
	repo repositories.LessonRepositoryInterface
}

func NewLessonService(repo repositories.LessonRepositoryInterface) LessonServiceInterface {
	return &LessonService{repo: repo}
}

// CreateLesson создает новый урок
func (s *LessonService) CreateLesson(ctx context.Context, lesson *models.Lesson) (*models.Lesson, error) {
	if err := s.repo.Create(ctx, lesson); err != nil {
		return nil, err
	}
	return lesson, nil
}

// GetLessonByID получает урок по ID
func (s *LessonService) GetLessonByID(ctx context.Context, id int) (*models.Lesson, error) {
	return s.repo.FindByID(ctx, id)
}

// GetAllLessons получает все уроки
func (s *LessonService) GetAllLessons(ctx context.Context, courseID int) ([]*models.Lesson, error) {
	return s.repo.FindByCourseID(ctx, courseID)
}

func (s *LessonService) DeleteLesson(ctx context.Context, id int) error {
	return s.repo.Delete(ctx, id)
}

// UpdateLesson обновляет урок
func (s *LessonService) UpdateLesson(ctx context.Context, lesson *models.Lesson) error {
	return s.repo.Update(ctx, lesson)
}