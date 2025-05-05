package usecase

import (
	"context"
	"errors"
	"log"
	"time"
	"database/sql"
	"gitlab.com/w0ikid/study-platform/internal/domain/models"
	"gitlab.com/w0ikid/study-platform/internal/domain/services"
	"gitlab.com/w0ikid/study-platform/internal/domain/services/pdfgen"
)

type CertificateUseCaseInterface interface {
	Generate(ctx context.Context, userID, courseID int) ([]byte, error)
	CreateCertificate(ctx context.Context, userID, courseID int) (*models.Certificate, error)
	GetCertificateByUserAndCourse(ctx context.Context, userID, courseID int) (*models.Certificate, error)
	GetCertificatesByUserID(ctx context.Context, userID int) ([]*models.Certificate, error)
}

type CertificateUseCase struct {
	certificateService services.CertificateServiceInterface
	enrollmentService  services.EnrollmentServiceInterface
	userService 	  services.UserServiceInterface
	courseService 	  services.CourseServiceInterface
}

func NewCertificateUseCase(certificateService services.CertificateServiceInterface, enrollmentService services.EnrollmentServiceInterface, userService services.UserServiceInterface, courseService services.CourseServiceInterface) *CertificateUseCase {
	return &CertificateUseCase{
		certificateService: certificateService,
		enrollmentService:  enrollmentService,
		userService:        userService,
		courseService:      courseService,
	}
}

func (uc *CertificateUseCase) CreateCertificate(ctx context.Context, userID, courseID int) (*models.Certificate, error) {
	enroll, err := uc.enrollmentService.GetEnrollmentByUserAndCourse(ctx, userID, courseID)
	
	if err != nil {
		return nil, err
	}
	
	
	if enroll == nil {
		return nil, errors.New("user is not enrolled in the course")
	}

	if enroll.Status != "completed" {
		return nil, errors.New("course not completed")
	}
	
	existing, err := uc.certificateService.GetCertificateByUserAndCourse(ctx, userID, courseID)
	if err == nil && existing != nil {
		return existing, nil
	}
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}
	
	certificate := &models.Certificate{
		UserID:   userID,
		CourseID: courseID,	
		IssuedAt: time.Now(),
	}
	return uc.certificateService.CreateCertificate(ctx, certificate)
} 

func (uc *CertificateUseCase) GetCertificateByUserAndCourse(ctx context.Context, userID, courseID int) (*models.Certificate, error) {
	certificate, err := uc.certificateService.GetCertificateByUserAndCourse(ctx, userID, courseID)
	if err != nil {
		return nil, err
	}
	return certificate, nil
}

func (uc *CertificateUseCase) GetCertificatesByUserID(ctx context.Context, userID int) ([]*models.Certificate, error) {
	certificates, err := uc.certificateService.GetCertificatesByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	return certificates, nil
}

func (uc *CertificateUseCase) Generate(ctx context.Context, userID, courseID int) ([]byte, error) {
	certificate, err := uc.CreateCertificate(ctx, userID, courseID)
	if err != nil {
		return nil, err
	}
	log.Print(1)
	user, err := uc.userService.GetUser(ctx, userID)
	if err != nil {
		return nil, err
	}
	log.Print(2)
	course, err := uc.courseService.GetCourse(ctx, courseID)
	if err != nil {
		return nil, err
	}
	log.Print(3)
	return pdfgen.GenerateCertificatePDF(user, course, certificate.IssuedAt)
}
