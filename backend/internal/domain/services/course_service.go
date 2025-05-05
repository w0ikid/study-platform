package services

import (
	"context"
	"gitlab.com/w0ikid/study-platform/internal/domain/models"
	"gitlab.com/w0ikid/study-platform/internal/domain/repositories"
)

type CourseServiceInterface interface {
	CreateCourse(ctx context.Context, course *models.Course) (*models.Course, error)
	GetCourse(ctx context.Context, id int) (*models.Course, error)
	GetAllCourses(ctx context.Context) ([]models.Course, error)
	UpdateCourse(ctx context.Context, course *models.Course) error
	DeleteCourse(ctx context.Context, id int) error
}

type CourseService struct {
	repo repositories.CourseRepositoryInterface
}

func NewCourseService(repo repositories.CourseRepositoryInterface) CourseServiceInterface {
	return &CourseService{repo: repo}
}

// CreateCourse создаёт новый курс
func (s *CourseService) CreateCourse(ctx context.Context, course *models.Course) (*models.Course, error) {
	if err := s.repo.Create(ctx, course); err != nil {
		return nil, err
	}
	return course, nil
}

// GetCourse возвращает курс по ID
func (s *CourseService) GetCourse(ctx context.Context, id int) (*models.Course, error) {
	return s.repo.FindByID(ctx, id)
}

// GetAllCourses возвращает все курсы
func (s *CourseService) GetAllCourses(ctx context.Context) ([]models.Course, error) {
	return s.repo.FindAll(ctx)
}

// UpdateCourse обновляет курс
func (s *CourseService) UpdateCourse(ctx context.Context, course *models.Course) error {
	return s.repo.Update(ctx, course)
}

// DeleteCourse удаляет курс
func (s *CourseService) DeleteCourse(ctx context.Context, id int) error {
	return s.repo.Delete(ctx, id)
}
