package repositories

import (
	"context"
	"gitlab.com/w0ikid/study-platform/internal/domain/models"
	"github.com/jackc/pgx/v5"
	"log"
)

type CertificateRepositoryInterface interface {
	CreateCertificate(ctx context.Context, certificate *models.Certificate) (*models.Certificate, error)
	GetCertificateByID(ctx context.Context, id int) (*models.Certificate, error)
	GetCertificatesByUserID(ctx context.Context, userID int) ([]*models.Certificate, error)
	GetCertificateByUserAndCourse(ctx context.Context, userID, courseID int) (*models.Certificate, error)

}

type CertificateRepository struct {
	db *pgx.Conn
}

func NewCertificateRepository(db *pgx.Conn) *CertificateRepository {
	return &CertificateRepository{db: db}
}

// CreateCertificate добавляет новый сертификат в базу данных
func (r *CertificateRepository) CreateCertificate(ctx context.Context, certificate *models.Certificate) (*models.Certificate, error) {
	query := `
		INSERT INTO certificates (user_id, course_id, issued_at)
		VALUES ($1, $2, $3)
		RETURNING id`
	err := r.db.QueryRow(ctx, query, certificate.UserID, certificate.CourseID, certificate.IssuedAt).
		Scan(&certificate.ID)
	if err != nil {
		log.Printf("Error creating certificate: %v", err)
		return nil, err
	}
	return certificate, nil
}

// GetCertificateByID получает сертификат по ID
func (r *CertificateRepository) GetCertificateByID(ctx context.Context, id int) (*models.Certificate, error) {
	var certificate models.Certificate
	query := `
		SELECT id, user_id, course_id, issued_at
		FROM certificates WHERE id = $1`
	err := r.db.QueryRow(ctx, query, id).
		Scan(&certificate.ID, &certificate.UserID, &certificate.CourseID, &certificate.IssuedAt)
	if err != nil {
		return nil, err
	}
	return &certificate, nil
}

// GetCertificatesByUserID получает все сертификаты пользователя по его ID
func (r *CertificateRepository) GetCertificatesByUserID(ctx context.Context, userID int) ([]*models.Certificate, error) {
	var certificates []*models.Certificate
	query := `
		SELECT id, user_id, course_id, issued_at
		FROM certificates WHERE user_id = $1`
	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var certificate models.Certificate
		if err := rows.Scan(&certificate.ID, &certificate.UserID, &certificate.CourseID, &certificate.IssuedAt); err != nil {
			return nil, err
		}
		certificates = append(certificates, &certificate)
	}

	return certificates, nil
}

// GetCertificateByUserAndCourse получает сертификат по ID пользователя и ID курса
func (r *CertificateRepository) GetCertificateByUserAndCourse(ctx context.Context, userID, courseID int) (*models.Certificate, error) {
	var certificate models.Certificate
	query := `
		SELECT id, user_id, course_id, issued_at
		FROM certificates WHERE user_id = $1 AND course_id = $2`
	err := r.db.QueryRow(ctx, query, userID, courseID).
		Scan(&certificate.ID, &certificate.UserID, &certificate.CourseID, &certificate.IssuedAt)
	if err != nil {
		return nil, err
	}
	return &certificate, nil
}