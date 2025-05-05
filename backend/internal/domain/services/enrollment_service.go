package services

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"gitlab.com/w0ikid/study-platform/internal/domain/models"
	"gitlab.com/w0ikid/study-platform/internal/domain/repositories"
)

type EnrollmentServiceInterface interface {
	CreateEnrollment(ctx context.Context, enrollment *models.Enrollment) (*models.Enrollment, error)
	GetEnrollmentByID(ctx context.Context, id int) (*models.Enrollment, error)
	IsUserEnrolled(ctx context.Context, userID, courseID int) (bool, error)
	MarkAsCompleted(ctx context.Context, userID, courseID int) error
	GetEnrollmentsByUser(ctx context.Context, userID int) ([]*models.Enrollment, error)
	GetEnrollmentsByCourse(ctx context.Context, courseID int) ([]*models.Enrollment, error)
	DeleteEnrollment(ctx context.Context, id int) error
	GetEnrollmentByUserAndCourse(ctx context.Context, userID, courseID int) (*models.Enrollment, error)
}
type EnrollmentService struct {
	repo repositories.EnrollmentRepositoryInterface
}

func NewEnrollmentService(repo repositories.EnrollmentRepositoryInterface) EnrollmentServiceInterface {
	return &EnrollmentService{repo: repo}
}

// CreateEnrollment создает новую запись о зачислении
func (s *EnrollmentService) CreateEnrollment(ctx context.Context, enrollment *models.Enrollment) (*models.Enrollment, error) {
	if err := s.repo.Create(ctx, enrollment); err != nil {
		return nil, err
	}
	return enrollment, nil
}

// GetEnrollmentByID получает запись о зачислении по ID
func (s *EnrollmentService) GetEnrollmentByID(ctx context.Context, id int) (*models.Enrollment, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *EnrollmentService) IsUserEnrolled(ctx context.Context, userID, courseID int) (bool, error) {
	enrollment, err := s.repo.FindByUserAndCourseID(ctx, userID, courseID)
	if err != nil {
		// допустимая ошибка no rows
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}
		return false, err
	}
	return enrollment != nil, nil
}

func (s *EnrollmentService) MarkAsCompleted(ctx context.Context, userID, courseID int) error {
	enrollmentID, err := s.repo.FindByUserAndCourseID(ctx, userID, courseID)
	if err != nil {
        return err
    }
	return s.repo.UpdateStatus(ctx, enrollmentID.ID, "completed")
}

func (s *EnrollmentService) GetEnrollmentByUserAndCourse(ctx context.Context, userID, courseID int) (*models.Enrollment, error) {
	return s.repo.FindByUserAndCourseID(ctx, userID, courseID)
}

func (s *EnrollmentService) GetEnrollmentsByUser(ctx context.Context, userID int) ([]*models.Enrollment, error) {
	return s.repo.FindByUserID(ctx, userID)
}

func (s *EnrollmentService) GetEnrollmentsByCourse(ctx context.Context, courseID int) ([]*models.Enrollment, error) {
	return s.repo.FindByCourseID(ctx, courseID)
}

func (s *EnrollmentService) DeleteEnrollment(ctx context.Context, id int) error {
	return s.repo.Delete(ctx, id)
}