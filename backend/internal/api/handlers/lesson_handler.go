package handlers

import (
	_"errors"
	_"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gitlab.com/w0ikid/study-platform/internal/domain/usecase"
	"gitlab.com/w0ikid/study-platform/internal/dto"
)


type LessonHandler struct {
	lessonUseCase *usecase.LessonUseCase
}

func NewLessonHandler(lessonUseCase *usecase.LessonUseCase) *LessonHandler {
	return &LessonHandler{
		lessonUseCase: lessonUseCase,
	}
}

// CreateLesson обрабатывает создание нового урока
func (h *LessonHandler) CreateLesson(c *gin.Context) {
	var request dto.CreateLessonRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	courseIDparam := c.Param("id")
	courseID, err := strconv.Atoi(courseIDparam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid course ID"})
		return
	}


	userID := c.GetInt("userID")

	ctx := c.Request.Context()

	input := usecase.CreateLessonInput{
		Title:     request.Title,
		Content:   request.Content,
		CourseID:  courseID,
		VideoURL:  request.VideoURL,
		UserID: userID,
	}

	lesson, err := h.lessonUseCase.CreateLesson(ctx, input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":          lesson.ID,
		"name":        lesson.Title,
	})
}

func (h *LessonHandler) GetLessonsByCourse(c *gin.Context) {
	userID := c.GetInt("userID")
	courseID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid course ID"})
		return
	}
	
	ctx := c.Request.Context()
	
	// h.lessonUseCase.GetLessonsForStudent(ctx, userID, courseID)
	lessons, err := h.lessonUseCase.GetLessonsForStudent(ctx, userID, courseID)
	
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "you are not enrolled for this course"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"lessons": lessons,
	})
}