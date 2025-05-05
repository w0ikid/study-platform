package usecase

import (
	"context"
	"errors"
	"fmt"
	"gitlab.com/w0ikid/study-platform/internal/domain/models"
	"gitlab.com/w0ikid/study-platform/internal/domain/services"
)

type LessonUseCase struct {
	lessonService services.LessonServiceInterface
	enrollment services.EnrollmentServiceInterface
	course services.CourseServiceInterface
}

func NewLessonUseCase(
	lessonService services.LessonServiceInterface, enrollment services.EnrollmentServiceInterface, course services.CourseServiceInterface,
) *LessonUseCase {
	return &LessonUseCase{lessonService: lessonService, enrollment: enrollment, course: course}
}

type CreateLessonInput struct {
	Title string
	Content string
	CourseID int
	VideoURL string
	UserID int
}

// CreateLesson создает новый урок
func (u *LessonUseCase) CreateLesson(ctx context.Context, input CreateLessonInput) (*models.Lesson, error) {
	isOwner, err := u.IsCourseOwnedByTeacher(ctx, input.CourseID, input.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to verify course ownership: %w", err)
	}
	if !isOwner {
		return nil, errors.New("user is not the owner of this course")
	}

	if input.Title == "" {
		return nil, errors.New("lesson title cannot be empty")
	}

	lesson := &models.Lesson{
		Title:     input.Title,
		Content:   input.Content,
		CourseID:  input.CourseID,
		VideoURL:  input.VideoURL,
	}
	

	
	lesson, err = u.lessonService.CreateLesson(ctx, lesson)
	if err != nil {
		return nil, err
	}
	// Дополнительная логика (например, отправка уведомлений)
	return lesson, nil
}

// GetLessonByID получает урок по ID
func (u *LessonUseCase) GetLessonByID(ctx context.Context, id int) (*models.Lesson, error) {
	// Дополнительная логика перед получением урока (например, проверка прав доступа)
	return u.lessonService.GetLessonByID(ctx, id)
}

// GetAllLessons получает все уроки
func (u *LessonUseCase) GetAllLessons(ctx context.Context, courseID int) ([]*models.Lesson, error) {
	// Дополнительная логика перед получением всех уроков (например, проверка прав доступа)
	return u.lessonService.GetAllLessons(ctx, courseID)
}

// DeleteLesson удаляет урок по ID
func (u *LessonUseCase) DeleteLesson(ctx context.Context, id int) error {
	// Дополнительная логика перед удалением урока (например, проверка прав доступа)
	return u.lessonService.DeleteLesson(ctx, id)
}

func (u *LessonUseCase) GetLessonsForStudent(ctx context.Context, userID int, courseID int) ([]*models.Lesson, error) {
	// enrolled, err := u.enrollment.IsUserEnrolled(ctx, userID, courseID)
	// if err != nil {
	// 	return nil, err
	// }
	// if !enrolled {
	// 	return nil, errors.New("access denied | you are not enrolled in this course")
	// }
	
	return u.lessonService.GetAllLessons(ctx, courseID)
}

func (uc *LessonUseCase) IsCourseOwnedByTeacher(ctx context.Context, courseID, teacherID int) (bool, error) {
    course, err := uc.course.GetCourse(ctx, courseID)
    if err != nil {
        return false, err
    }

    return course.TeacherID == teacherID, nil
}
