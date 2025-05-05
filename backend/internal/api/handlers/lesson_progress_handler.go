package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gitlab.com/w0ikid/study-platform/internal/domain/usecase"
)

type LessonProgressHandler struct {
	lessonProgressUseCase *usecase.LessonProgressUseCase
}

func NewLessonProgressHandler(lessonProgressUseCase *usecase.LessonProgressUseCase) *LessonProgressHandler {
	return &LessonProgressHandler{lessonProgressUseCase: lessonProgressUseCase}
}

func (h *LessonProgressHandler) CompleteLesson(c *gin.Context) {
	userID := c.GetInt("userID")
	courseID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid course_id"})
        return
	}

	lessonID, err := strconv.Atoi(c.Param("lesson_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid lesson_id"})
        return
	}

	ctx := c.Request.Context()

	err = h.lessonProgressUseCase.MarkLessonCompleted(ctx, userID, lessonID, courseID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Lesson marked as completed and XP awarded"})
}

func (h *LessonProgressHandler) GetCourseProgress(c *gin.Context) {
	userID := c.GetInt("userID")
	courseID, err := strconv.Atoi(c.Param("id")) 

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid courseID"})
        return
	}

	ctx := c.Request.Context()

	progress, err := h.lessonProgressUseCase.GetCourseProgress(ctx, userID, courseID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid progress"})
        return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"progress": progress,
	})
}