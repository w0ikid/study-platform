package handlers

import (
	"net/http"
	"strconv"
	"github.com/gin-gonic/gin"
	"gitlab.com/w0ikid/study-platform/internal/domain/usecase"
)

type CertificateHandler struct {
	certificateUseCase *usecase.CertificateUseCase
}

func NewCertificateHandler(certificateUseCase *usecase.CertificateUseCase) *CertificateHandler {
	return &CertificateHandler{
		certificateUseCase: certificateUseCase,
	}
}

func (h *CertificateHandler) GenerateCertificate(c *gin.Context) {
	userID := c.GetInt("userID")
	courseID, err := strconv.Atoi(c.Param("course_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid course_id"})
		return
	}

	ctx := c.Request.Context()

	pdfData, err := h.certificateUseCase.Generate(ctx, userID, courseID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate certificate"})
		return
	}

	c.Header("Content-Type", "application/pdf")
	c.Header("Content-Disposition", "attachment; filename=certificate.pdf")
	c.Data(http.StatusOK, "application/pdf", pdfData)
	c.JSON(http.StatusOK, gin.H{"message": "Certificate generated successfully"})
}