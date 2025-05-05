package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gitlab.com/w0ikid/study-platform/internal/domain/usecase"
	"gitlab.com/w0ikid/study-platform/internal/dto"
)

type CourseHandler struct {
	courseUseCase *usecase.CourseUseCase
}

func NewCourseHandler(courseUseCase *usecase.CourseUseCase) *CourseHandler {
	return &CourseHandler{
		courseUseCase: courseUseCase,
	}
}

// CreateCourse обрабатывает создание нового курса
func (h *CourseHandler) CreateCourse(c *gin.Context) {
	var request dto.CreateCourseRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	teacherID := c.GetInt(("userID"))

	ctx := c.Request.Context()

	input := usecase.CreateCourseInput{
		Name:        request.Name,
		Description: request.Description,
		TeacherID:   teacherID,
		Status: 	 "active",
		ImageURL: 	 request.ImageUrl,
	}
	
	course, err := h.courseUseCase.CreateCourse(ctx, input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":          course.ID,
		"name":        course.Name,
		"description": course.Description,
		"teacher_id":  course.TeacherID,
		"image_url":   course.ImageUrl,
	})
}

// GetCourse обрабатывает получение курса по ID
func (h *CourseHandler) GetCourse(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid course ID"})
		return
	}

	course, err := h.courseUseCase.GetCourseByID(c.Request.Context(), int(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Course not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":          course.ID,
		"name":        course.Name,
		"description": course.Description,
		"teacher_id":  course.TeacherID,
	})
}

// GetAllCourses обрабатывает получение всех курсов
func (h *CourseHandler) GetAllCourses(c *gin.Context) {
	courses, err := h.courseUseCase.GetAllCourses(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"courses": courses})
}

func (h *CourseHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid course ID"})
		return 
	}

	ctx := c.Request.Context()

	err = h.courseUseCase.DeleteCourse(ctx, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Course not found"})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}