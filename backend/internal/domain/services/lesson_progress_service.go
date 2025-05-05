package services

import (
	"context"
	"fmt"
	"time"

	"gitlab.com/w0ikid/study-platform/internal/domain/models"
	"gitlab.com/w0ikid/study-platform/internal/domain/repositories"
)

type LessonProgressServiceInterface interface {
    MarkLessonCompleted(ctx context.Context, userID, lessonID, courseID int) error
    GetProgressByCourse(ctx context.Context, userID, courseID int) ([]*models.LessonProgress, error)
    GetProgressByLesson(ctx context.Context, userID, lessonID int) (*models.LessonProgress, error)
}

type LessonProgressService struct {
    repo repositories.LessonProgressRepositoryInterface
}

func NewLessonProgressService(repo repositories.LessonProgressRepositoryInterface) *LessonProgressService {
    return &LessonProgressService{repo: repo}
}

func (s *LessonProgressService) MarkLessonCompleted(ctx context.Context, userID, lessonID, courseID int) error {
    progress, _ := s.repo.FindByUserAndLesson(ctx, userID, lessonID)
    // if err != nil && err.Error() != "failed to find lesson progress: record not found" {
    //     return err
    // }
    fmt.Println(progress)
    if progress == nil {
        now := time.Now()
        progress = &models.LessonProgress{
            UserID:      userID,
            LessonID:    lessonID,
            CourseID:    courseID,
            IsCompleted: true,
            CompletedAt: now,
        }
        return s.repo.Create(ctx, progress)
    }
    now := time.Now()
    progress.IsCompleted = true
    progress.CompletedAt = now
    return s.repo.Update(ctx, progress)
}

func (s *LessonProgressService) GetProgressByCourse(ctx context.Context, userID, courseID int) ([]*models.LessonProgress, error) {
    return s.repo.FindByUserAndCourse(ctx, userID, courseID)
}

func (s *LessonProgressService) GetProgressByLesson(ctx context.Context, userID, lessonID int) (*models.LessonProgress, error) {
    return s.repo.FindByUserAndLesson(ctx, userID, lessonID)
}