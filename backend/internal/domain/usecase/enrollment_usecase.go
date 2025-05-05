package usecase

import (
	"context"
	"errors"
	_"log"

	"gitlab.com/w0ikid/study-platform/internal/domain/models"
	"gitlab.com/w0ikid/study-platform/internal/domain/services"
)

type EnrollmentUseCaseInterface interface {
	EnrollStudent(ctx context.Context, userID, courseID int) error
	GetEnrollmentByID(ctx context.Context, id int) (*models.Enrollment, error)
	IsUserEnrolled(ctx context.Context, userID, courseID int) (bool, error)
	GetStudentEnrollments(ctx context.Context, userID int) ([]*models.Enrollment, error)
	GetCourseEnrollments(ctx context.Context, courseID int) ([]*models.Enrollment, error)
	DeleteEnrollment(ctx context.Context, id int) error
}

type EnrollmentUseCase struct {
	enrollmentService services.EnrollmentServiceInterface
	courseService     services.CourseServiceInterface
}
func NewEnrollmentUseCase(
	enrollmentService services.EnrollmentServiceInterface,
	courseService services.CourseServiceInterface,
) *EnrollmentUseCase {
	return &EnrollmentUseCase{
		enrollmentService: enrollmentService,
		courseService:     courseService,
	}
}
func (u *EnrollmentUseCase) EnrollStudent(ctx context.Context, userID, courseID int) error {
	// check if course exists
	course, err := u.courseService.GetCourse(ctx, courseID)
	if err != nil {
		return err
	}
	if course == nil {
		return errors.New("course not found")
	}
	
	if course.TeacherID == userID {
		return errors.New("teacher cannot enroll in their own course")
	}
	
	// check if already enrolled
	enrolled, err := u.enrollmentService.IsUserEnrolled(ctx, userID, courseID)
	if err != nil {
		return err
	}
	if enrolled {
		return errors.New("student already enrolled in this course")
	} 

	
	// Create enrollment
	enrollment := &models.Enrollment{
		UserID:   userID,
		CourseID: courseID,
		Status:   "active",
	}

	_, err = u.enrollmentService.CreateEnrollment(ctx, enrollment)
	return err
}

// GetEnrollmentByID gets an enrollment record by ID
func (u *EnrollmentUseCase) GetEnrollmentByID(ctx context.Context, id int) (*models.Enrollment, error) {
	return u.enrollmentService.GetEnrollmentByID(ctx, id)
}

func (u *EnrollmentUseCase) IsUserEnrolled(ctx context.Context, userID, courseID int) (bool, error) {
	return u.enrollmentService.IsUserEnrolled(ctx, userID, courseID)
}

func (u *EnrollmentUseCase) GetStudentEnrollments(ctx context.Context, userID int) ([]*models.Enrollment, error) {
	return u.enrollmentService.GetEnrollmentsByUser(ctx, userID)
}

func (u *EnrollmentUseCase) GetCourseEnrollments(ctx context.Context, courseID int) ([]*models.Enrollment, error) {
	return u.enrollmentService.GetEnrollmentsByCourse(ctx, courseID)
}

func (u *EnrollmentUseCase) DeleteEnrollment(ctx context.Context, id int) error {
	return u.enrollmentService.DeleteEnrollment(ctx, id)
}