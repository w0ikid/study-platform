package repositories

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"gitlab.com/w0ikid/study-platform/internal/domain/models"
)

type LessonRepositoryInterface interface {
	Create(ctx context.Context, lesson *models.Lesson) error
	FindByID(ctx context.Context, id int) (*models.Lesson, error)
	FindByCourseID(ctx context.Context, courseID int) ([]*models.Lesson, error)
	Update(ctx context.Context, lesson *models.Lesson) error
	Delete(ctx context.Context, id int) error
}
type LessonRepository struct {
	db *pgx.Conn
}

func NewLessonRepository(db *pgx.Conn) *LessonRepository {
	return &LessonRepository{db: db}
}

// Create inserts a new lesson into the database and returns the created lesson
func (r *LessonRepository) Create(ctx context.Context, lesson *models.Lesson) error {
	query := `
		INSERT INTO lessons (course_id, title, content, video_url)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, updated_at`
	err := r.db.QueryRow(ctx, query, lesson.CourseID, lesson.Title, lesson.Content, lesson.VideoURL).
		Scan(&lesson.ID, &lesson.CreatedAt, &lesson.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to create lesson: %w", err)
	}
	return nil
}

// FindByID retrieves a lesson by its ID
func (r *LessonRepository) FindByID(ctx context.Context, id int) (*models.Lesson, error) {
	query := `
		SELECT id, course_id, title, content, video_url, created_at, updated_at
		FROM lessons
		WHERE id = $1`
	var lesson models.Lesson
	err := r.db.QueryRow(ctx, query, id).Scan(
		&lesson.ID,
		&lesson.CourseID,
		&lesson.Title,
		&lesson.Content,
		&lesson.VideoURL,
		&lesson.CreatedAt,
		&lesson.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to find lesson by id: %w", err)
	}
	return &lesson, nil
}

// FindByCourseID retrieves all lessons for a given course ID
func (r *LessonRepository) FindByCourseID(ctx context.Context, courseID int) ([]*models.Lesson, error) {
	query := `
		SELECT id, course_id, title, content, video_url, created_at, updated_at
		FROM lessons
		WHERE course_id = $1`
	rows, err := r.db.Query(ctx, query, courseID)
	if err != nil {
		return nil, fmt.Errorf("failed to find lessons by course id: %w", err)
	}
	defer rows.Close()

	var lessons []*models.Lesson
	for rows.Next() {
		var lesson models.Lesson
		if err := rows.Scan(
			&lesson.ID,
			&lesson.CourseID,
			&lesson.Title,
			&lesson.Content,
			&lesson.VideoURL,
			&lesson.CreatedAt,
			&lesson.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan lesson: %w", err)
		}
		lessons = append(lessons, &lesson)
	}
	return lessons, nil
}

// Update updates an existing lesson in the database
func (r *LessonRepository) Update(ctx context.Context, lesson *models.Lesson) error {
	query := `
		UPDATE lessons
		SET title = $1, content = $2, video_url = $3, updated_at = NOW()
		WHERE id = $4`
	_, err := r.db.Exec(ctx, query, lesson.Title, lesson.Content, lesson.VideoURL, lesson.ID)
	if err != nil {
		return fmt.Errorf("failed to update lesson: %w", err)
	}	
	return nil
}

// Delete removes a lesson from the database
func (r *LessonRepository) Delete(ctx context.Context, id int) error {
	query := `
		DELETE FROM lessons
		WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete lesson: %w", err)
	}
	return nil
}


// LessonProgress -----------------------

