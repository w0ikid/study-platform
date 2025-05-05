package usecase

import (
	"context"
	"gitlab.com/w0ikid/study-platform/internal/domain/models"
	"gitlab.com/w0ikid/study-platform/internal/domain/services"	
)

type CourseUseCaseInterface interface {
	CreateCourse(ctx context.Context, input CreateCourseInput) (*models.Course, error)
	GetCourseByID(ctx context.Context, id int) (*models.Course, error)
	GetAllCourses(ctx context.Context) ([]models.Course, error)
	UpdateCourse(ctx context.Context, course *models.Course) error
	DeleteCourse(ctx context.Context, id int) error
}

type CourseUseCase struct {
	courseService services.CourseServiceInterface
}

func NewCourseUseCase(courseService services.CourseServiceInterface) *CourseUseCase {
	return &CourseUseCase{courseService: courseService}
}

type CreateCourseInput struct {
	Name        string
	Description	string
	TeacherID   int
	Status      string
	ImageURL 	string
}

// CreateCourse создает новый курс
func (u *CourseUseCase) CreateCourse(ctx context.Context, input CreateCourseInput) (*models.Course, error) {


	course := &models.Course{
		Name:        input.Name,
		Description: input.Description,
		TeacherID:   input.TeacherID,
		Status:      input.Status,
		ImageUrl:    input.ImageURL,
	}
	
	// Вызов сервиса для создания курса
	course, err := u.courseService.CreateCourse(ctx, course)
	if err != nil {
		return nil, err
	}
	return course, nil
}

// GetCourseByID получает курс по ID
func (u *CourseUseCase) GetCourseByID(ctx context.Context, id int) (*models.Course, error) {
	// Дополнительная логика перед получением курса (например, проверка прав доступа)
	return u.courseService.GetCourse(ctx, id)
}

// GetAllCourses получает все курсы
func (u *CourseUseCase) GetAllCourses(ctx context.Context) ([]models.Course, error) {
	// Дополнительная логика перед получением всех курсов (например, фильтрация)
	return u.courseService.GetAllCourses(ctx)
}

// UpdateCourse обновляет курс
func (u *CourseUseCase) UpdateCourse(ctx context.Context, course *models.Course) error {
	// Дополнительная логика перед обновлением курса (например, проверка прав доступа)
	return u.courseService.UpdateCourse(ctx, course)
}

// DeleteCourse удаляет курс
func (u *CourseUseCase) DeleteCourse(ctx context.Context, id int) error {
	// Дополнительная логика перед удалением курса (например, проверка прав доступа)
	return u.courseService.DeleteCourse(ctx, id)
}