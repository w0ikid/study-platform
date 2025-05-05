package repositories

import (
	"context"
	"fmt"
    "errors"
	"github.com/jackc/pgx/v5"
    "gitlab.com/w0ikid/study-platform/internal/domain/models"
)

type LessonProgressRepositoryInterface interface {
    Create(ctx context.Context, progress *models.LessonProgress) error
    FindByUserAndLesson(ctx context.Context, userID, lessonID int) (*models.LessonProgress, error)
    Update(ctx context.Context, progress *models.LessonProgress) error
    FindByUserAndCourse(ctx context.Context, userID, courseID int) ([]*models.LessonProgress, error)
}

type LessonProgressRepository struct {
    db *pgx.Conn
}

func NewLessonProgressRepository(db *pgx.Conn) *LessonProgressRepository {
    return &LessonProgressRepository{db: db}
}

func (r *LessonProgressRepository) Create(ctx context.Context, progress *models.LessonProgress) error {
    query := `
        INSERT INTO lesson_progress (user_id, lesson_id, course_id, is_completed, completed_at)
        VALUES ($1, $2, $3, $4, $5)
        RETURNING id, created_at, updated_at`
    err := r.db.QueryRow(ctx, query,
        progress.UserID,
        progress.LessonID,
        progress.CourseID,
        progress.IsCompleted,
        progress.CompletedAt,
    ).Scan(&progress.ID, &progress.CreatedAt, &progress.UpdatedAt)
    if err != nil {
        return fmt.Errorf("failed to create lesson progress: %w", err)
    }
    return nil
}

func (r *LessonProgressRepository) FindByUserAndLesson(ctx context.Context, userID, lessonID int) (*models.LessonProgress, error) {
    query := `
        SELECT id, user_id, lesson_id, course_id, is_completed, completed_at, created_at, updated_at
        FROM lesson_progress
        WHERE user_id = $1 AND lesson_id = $2`
    
    var progress models.LessonProgress
    err := r.db.QueryRow(ctx, query, userID, lessonID).Scan(
        &progress.ID,
        &progress.UserID,
        &progress.LessonID,
        &progress.CourseID,
        &progress.IsCompleted,
        &progress.CompletedAt,
        &progress.CreatedAt,
        &progress.UpdatedAt,
    )

    if err != nil {
        if errors.Is(err, pgx.ErrNoRows) {
            return nil, nil // запись не найдена - это не ошибка :3
        }
        return nil, fmt.Errorf("failed to find lesson progress: %w", err)
    }

    return &progress, nil
}


func (r *LessonProgressRepository) Update(ctx context.Context, progress *models.LessonProgress) error {
    query := `
        UPDATE lesson_progress
        SET is_completed = $1, completed_at = $2, updated_at = NOW()
        WHERE id = $3`
    _, err := r.db.Exec(ctx, query, progress.IsCompleted, progress.CompletedAt, progress.ID)
    if err != nil {
        return fmt.Errorf("failed to update lesson progress: %w", err)
    }
    return nil
}

func (r *LessonProgressRepository) FindByUserAndCourse(ctx context.Context, userID, courseID int) ([]*models.LessonProgress, error) {
    query := `
        SELECT id, user_id, lesson_id, course_id, is_completed, completed_at, created_at, updated_at
        FROM lesson_progress
        WHERE user_id = $1 AND course_id = $2`
    rows, err := r.db.Query(ctx, query, userID, courseID)
    if err != nil {
        return nil, fmt.Errorf("failed to find lesson progress by user and course: %w", err)
    }
    defer rows.Close()

    var progresses []*models.LessonProgress
    for rows.Next() {
        var progress models.LessonProgress
        if err := rows.Scan(
            &progress.ID,
            &progress.UserID,
            &progress.LessonID,
            &progress.CourseID,
            &progress.IsCompleted,
            &progress.CompletedAt,
            &progress.CreatedAt,
            &progress.UpdatedAt,
        ); err != nil {
            return nil, fmt.Errorf("failed to scan lesson progress: %w", err)
        }
        progresses = append(progresses, &progress)
    }
    return progresses, nil
}