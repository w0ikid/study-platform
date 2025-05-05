package repositories

import (
	"context"
	"fmt"
	"errors"
	"github.com/jackc/pgx/v5"
	"gitlab.com/w0ikid/study-platform/internal/domain/models"
)

type EnrollmentRepositoryInterface interface {
	Create(ctx context.Context, enrollment *models.Enrollment) error
	FindByID(ctx context.Context, id int) (*models.Enrollment, error)
	FindByUserAndCourseID(ctx context.Context, userID, courseID int) (*models.Enrollment, error)
	FindByUserID(ctx context.Context, userID int) ([]*models.Enrollment, error)
	FindByCourseID(ctx context.Context, courseID int) ([]*models.Enrollment, error)
	UpdateStatus(ctx context.Context, id int, status string) error
	Delete(ctx context.Context, id int) error
}

type EnrollmentRepository struct {
	db *pgx.Conn
}

func NewEnrollmentRepository(db *pgx.Conn) *EnrollmentRepository {
	return &EnrollmentRepository{db: db}
}

// Create добавляет новую запись о записи на курс
func (r *EnrollmentRepository) Create(ctx context.Context, enrollment *models.Enrollment) error {
	query := `INSERT INTO enrollments (user_id, course_id, created_at, updated_at)
		VALUES ($1, $2, NOW(), NOW()) RETURNING id`
	err := r.db.QueryRow(ctx, query,
		enrollment.UserID, enrollment.CourseID).
		Scan(&enrollment.ID)
	if err != nil {
		return fmt.Errorf("failed to create enrollment: %w", err)
	}
	return nil
}

// FindByID ищет запись о записи на курс по ID
func (r *EnrollmentRepository) FindByID(ctx context.Context, id int) (*models.Enrollment, error) {
	var enrollment models.Enrollment
	query := `SELECT id, user_id, course_id, status, created_at, updated_at FROM enrollments WHERE id = $1`

	err := r.db.QueryRow(ctx, query, id).Scan(&enrollment.ID, &enrollment.UserID, &enrollment.CourseID, &enrollment.Status, &enrollment.CreatedAt, &enrollment.UpdatedAt)


	if err != nil {
		return nil, fmt.Errorf("enrollment not found: %w", err)
	}

	return &enrollment, nil
}

func (r *EnrollmentRepository) FindByUserAndCourseID(ctx context.Context, userID, courseID int) (*models.Enrollment, error) {
	var enrollment models.Enrollment
	query := `SELECT id, user_id, course_id, status, created_at, updated_at FROM enrollments WHERE user_id = $1 AND course_id = $2`

	err := r.db.QueryRow(ctx, query, userID, courseID).Scan(&enrollment.ID, &enrollment.UserID, &enrollment.CourseID, &enrollment.Status, &enrollment.CreatedAt, &enrollment.UpdatedAt)
	
	if err != nil {
        if errors.Is(err, pgx.ErrNoRows) {
            return nil, fmt.Errorf("enrollment not found: %w", err)
        }
        return nil, fmt.Errorf("error finding enrollment: %w", err)
    }

	return &enrollment, nil
}

// find by userID
func (r *EnrollmentRepository) FindByUserID(ctx context.Context, userID int) ([]*models.Enrollment, error) {
	query := `SELECT id, user_id, course_id, status, created_at, updated_at FROM enrollments WHERE user_id = $1`

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("error querying enrollments: %w", err)
	}
	defer rows.Close()

	var enrollments []*models.Enrollment
	for rows.Next() {
		var enrollment models.Enrollment
		err := rows.Scan(&enrollment.ID, &enrollment.UserID, &enrollment.CourseID, &enrollment.Status, &enrollment.CreatedAt, &enrollment.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("error scanning enrollment: %w", err)
		}
		enrollments = append(enrollments, &enrollment)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return enrollments, nil
}


// find by courseID
func (r *EnrollmentRepository) FindByCourseID(ctx context.Context, courseID int) ([]*models.Enrollment, error) {
	query := `SELECT id, user_id, course_id, status, created_at, updated_at FROM enrollments WHERE course_id = $1`

	rows, err := r.db.Query(ctx, query, courseID)
	if err != nil {
		return nil, fmt.Errorf("error querying enrollments: %w", err)
	}
	defer rows.Close()

	var enrollments []*models.Enrollment
	for rows.Next() {
		var enrollment models.Enrollment
		err := rows.Scan(&enrollment.ID, &enrollment.UserID, &enrollment.CourseID, &enrollment.Status, &enrollment.CreatedAt, &enrollment.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("error scanning enrollment: %w", err)
		}
		enrollments = append(enrollments, &enrollment)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return enrollments, nil
}


// UpdateStatus обновляет только статус записи на курс
func (r *EnrollmentRepository) UpdateStatus(ctx context.Context, id int, status string) error {
    query := `UPDATE enrollments SET status = $1, updated_at = NOW() WHERE id = $2`
    
    commandTag, err := r.db.Exec(ctx, query, status, id)
    if err != nil {
        return fmt.Errorf("failed to update enrollment status: %w", err)
    }
    
    if commandTag.RowsAffected() == 0 {
        return fmt.Errorf("enrollment with id %d not found", id)
    }
    
    return nil
}

func (r *EnrollmentRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM enrollments WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}