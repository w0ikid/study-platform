package services

import (
	"context"
	"gitlab.com/w0ikid/study-platform/internal/domain/models"
	"gitlab.com/w0ikid/study-platform/internal/domain/repositories"
)

type CertificateServiceInterface interface {
	CreateCertificate(ctx context.Context, certificate *models.Certificate) (*models.Certificate, error)
	GetCertificateByID(ctx context.Context, id int) (*models.Certificate, error)
	GetCertificatesByUserID(ctx context.Context, userID int) ([]*models.Certificate, error)
	GetCertificateByUserAndCourse(ctx context.Context, userID, courseID int) (*models.Certificate, error)
}

type CertificateService struct {
	certificateRepo repositories.CertificateRepositoryInterface
}

func NewCertificateService(certificateRepo repositories.CertificateRepositoryInterface) *CertificateService {
	return &CertificateService{
		certificateRepo: certificateRepo,
	}
}

// CreateCertificate создает новый сертификат
func (s *CertificateService) CreateCertificate(ctx context.Context, certificate *models.Certificate) (*models.Certificate, error) {
	return s.certificateRepo.CreateCertificate(ctx, certificate)
}

// GetCertificateByID получает сертификат по ID
func (s *CertificateService) GetCertificateByID(ctx context.Context, id int) (*models.Certificate, error) {
	return s.certificateRepo.GetCertificateByID(ctx, id)
}

// GetCertificatesByUserID получает все сертификаты пользователя по его ID
func (s *CertificateService) GetCertificatesByUserID(ctx context.Context, userID int) ([]*models.Certificate, error) {
	return s.certificateRepo.GetCertificatesByUserID(ctx, userID)
}

// GetCertificateByUserAndCourse получает сертификат по ID пользователя и ID курса
func (s *CertificateService) GetCertificateByUserAndCourse(ctx context.Context, userID, courseID int) (*models.Certificate, error) {
	return s.certificateRepo.GetCertificateByUserAndCourse(ctx, userID, courseID)
}

