package handlers

import (
    "net/http"
    "strconv"

    "github.com/gin-gonic/gin"
    "gitlab.com/w0ikid/study-platform/internal/domain/usecase"
)

type EnrollmentHandler struct {
    enrollmentUseCase *usecase.EnrollmentUseCase
}

func NewEnrollmentHandler(enrollmentUseCase *usecase.EnrollmentUseCase) *EnrollmentHandler {
    return &EnrollmentHandler{
        enrollmentUseCase: enrollmentUseCase,
    }
}

// CreateEnrollment обрабатывает создание новой записи о зачислении
func (h *EnrollmentHandler) CreateEnrollment(c *gin.Context) {
    ctx := c.Request.Context()

    userID := c.GetInt("userID")
    
    courseIDparam := c.Param("id")
    courseID, err := strconv.Atoi(courseIDparam)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid course ID"})
        return
    }

    err = h.enrollmentUseCase.EnrollStudent(ctx, userID, courseID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusCreated, gin.H{"message": "Enrollment created successfully"})
}

// GetEnrollment обрабатывает получение записи о зачислении по ID
func (h *EnrollmentHandler) GetEnrollment(c *gin.Context) {
    idStr := c.Param("id")
    id, err := strconv.Atoi(idStr)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid enrollment ID"})
        return
    }

    ctx := c.Request.Context()

    enrollment, err := h.enrollmentUseCase.GetEnrollmentByID(ctx, id)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, enrollment)
}

// GetAllEnrollment обрабатывает получение всех записей о зачислении пользователя
func (h *EnrollmentHandler) GetAllEnrollment(c *gin.Context) {
    ctx := c.Request.Context()

    userID := c.GetInt("userID") // Предполагаем, что userID берется из контекста (authMiddleware)

    enrollments, err := h.enrollmentUseCase.GetStudentEnrollments(ctx, userID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, enrollments)
}

// Delete обрабатывает удаление записи о зачислении по ID
func (h *EnrollmentHandler) Delete(c *gin.Context) {
    idStr := c.Param("id")
    id, err := strconv.Atoi(idStr)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid enrollment ID"})
        return
    }

    ctx := c.Request.Context()

    err = h.enrollmentUseCase.DeleteEnrollment(ctx, id)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Enrollment deleted successfully"})
}