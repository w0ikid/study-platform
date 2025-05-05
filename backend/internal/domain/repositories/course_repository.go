package repositories

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"gitlab.com/w0ikid/study-platform/internal/domain/models"
)

type CourseRepositoryInterface interface {
	Create(ctx context.Context, course *models.Course) error
	FindByID(ctx context.Context, id int) (*models.Course, error)
	FindAll(ctx context.Context) ([]models.Course, error)
	Update(ctx context.Context, course *models.Course) error
	Delete(ctx context.Context, id int) error
}

type CourseRepository struct {
	db *pgx.Conn
}

func NewCourseRepository(db *pgx.Conn) *CourseRepository {
	return &CourseRepository{db: db}
}

// Create добавляет новый курс
func (r *CourseRepository) Create(ctx context.Context, course *models.Course) error {
	query := `
		INSERT INTO courses (name, description, image_url ,teacher_id, status)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, updated_at`
	err := r.db.QueryRow(ctx, query, course.Name, course.Description, course.ImageUrl ,course.TeacherID, course.Status).
		Scan(&course.ID, &course.CreatedAt, &course.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to create course: %w", err)
	}
	return nil
}

// FindByID ищет курс по ID
func (r *CourseRepository) FindByID(ctx context.Context, id int) (*models.Course, error) {
	var course models.Course
	query := `
		SELECT id, name, description, image_url ,teacher_id, status, created_at, updated_at
		FROM courses WHERE id = $1`
	err := r.db.QueryRow(ctx, query, id).
		Scan(&course.ID, &course.Name, &course.Description, &course.ImageUrl, &course.TeacherID, &course.Status, &course.CreatedAt, &course.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("course not found: %w", err)
	}
	return &course, nil
}

// FindAll возвращает список всех курсов
func (r *CourseRepository) FindAll(ctx context.Context) ([]models.Course, error) {
	query := `
		SELECT id, name, description, image_url,teacher_id, status, created_at, updated_at
		FROM courses`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch courses: %w", err)
	}
	defer rows.Close()

	var courses []models.Course
	for rows.Next() {
		var course models.Course
		if err := rows.Scan(&course.ID, &course.Name, &course.Description, &course.ImageUrl ,&course.TeacherID, &course.Status, &course.CreatedAt, &course.UpdatedAt); err != nil {
			return nil, fmt.Errorf("error scanning course: %w", err)
		}
		courses = append(courses, course)
	}

	return courses, nil
}

// Update обновляет курс
func (r *CourseRepository) Update(ctx context.Context, course *models.Course) error {
	query := `
		UPDATE courses 
		SET name = $1, description = $2, status = $3, updated_at = NOW()
		WHERE id = $4`
	_, err := r.db.Exec(ctx, query, course.Name, course.Description, course.Status, course.ID)
	if err != nil {
		return fmt.Errorf("failed to update course: %w", err)
	}
	return nil
}

// Delete удаляет курс
func (r *CourseRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM courses WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete course: %w", err)
	}
	return nil
}
