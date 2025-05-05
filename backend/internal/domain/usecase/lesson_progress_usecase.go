package usecase

import (
	"context"
	"errors"
	_ "log"

	"gitlab.com/w0ikid/study-platform/internal/domain/services"
)

type LessonProgressUseCaseInterface interface {
    MarkLessonCompleted(ctx context.Context, userID, lessonID, courseID int) error
    GetCourseProgress(ctx context.Context, userID, courseID int) (float64, error)
}

type LessonProgressUseCase struct {
    lessonProgressService services.LessonProgressServiceInterface
    lessonService         services.LessonServiceInterface
    enrollmentService     services.EnrollmentServiceInterface
    courseService         services.CourseServiceInterface
    userService           services.UserServiceInterface
}

func NewLessonProgressUseCase(
    lessonProgressService services.LessonProgressServiceInterface,
    lessonService services.LessonServiceInterface,
    enrollmentService services.EnrollmentServiceInterface,
    courseService services.CourseServiceInterface,
    userService services.UserServiceInterface,
) *LessonProgressUseCase {
    return &LessonProgressUseCase{
        lessonProgressService: lessonProgressService,
        lessonService:         lessonService,
        enrollmentService:     enrollmentService,
        courseService: courseService,
        userService: userService,
    }
}

func calculateLevel(xp int) int {
	return (xp / 100) + 1
}


func (uc *LessonProgressUseCase) MarkLessonCompleted(ctx context.Context, userID, lessonID, courseID int) error {
    // Проверяем, что пользователь записан на курс
    // enrolled, err := uc.enrollmentService.IsUserEnrolled(ctx, userID, courseID)
    // if err != nil {
    //     return err
    // }
    // if !enrolled {
    //     return errors.New("user is not enrolled in the course")
    // }
    
    // Проверяем, что урок принадлежит курсу
    lesson, err := uc.lessonService.GetLessonByID(ctx, lessonID)
    if err != nil {
        return err
    }
    
    if lesson.CourseID != courseID {
        return errors.New("lesson does not belong to the course")
    }

    lessonProg, err := uc.lessonProgressService.GetProgressByLesson(ctx, userID, lessonID)
    if err != nil {
        return err
    }

    if lessonProg != nil && lessonProg.IsCompleted {
        return errors.New("lesson is already completed")
    }

    const xpPerLesson = 10
    user, err := uc.userService.GetUser(ctx, userID)
    if err != nil {
        return errors.New("user not found in lesson_progress")
    }
    
    // +xp and level
    
    user.Xp = user.Xp + xpPerLesson
    user.Level = calculateLevel(user.Xp)
    
    uc.userService.UpdateXpAndLevel(ctx, user)
    // Отмечаем урок как завершенный
    
    return uc.lessonProgressService.MarkLessonCompleted(ctx, userID, lessonID, courseID)
}

func (uc *LessonProgressUseCase) GetCourseProgress(ctx context.Context, userID, courseID int) (float64, error) {
    // Проверяем, что пользователь записан на курс
    // enrolled, err := uc.enrollmentService.IsUserEnrolled(ctx, userID, courseID)
    // if err != nil {
    //     return 0, err
    // }
    // if !enrolled {
    //     return 0, errors.New("user is not enrolled in the course")
    // }

    // Получаем все уроки курса
    lessons, err := uc.lessonService.GetAllLessons(ctx, courseID)
    if err != nil {
        return 0, err
    }
    if len(lessons) == 0 {
        return 0, nil
    }

    // Получаем прогресс пользователя по курсу
    progresses, err := uc.lessonProgressService.GetProgressByCourse(ctx, userID, courseID)
    if err != nil {
        return 0, err
    }

    // Считаем процент завершения
    completedCount := 0
    for _, progress := range progresses {
        if progress.IsCompleted {
            completedCount++
        }
    }

    if completedCount == len(lessons) {
        uc.enrollmentService.MarkAsCompleted(ctx, userID, courseID)
    }

    return float64(completedCount) / float64(len(lessons)) * 100, nil
}