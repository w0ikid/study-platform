package models

import (
	"time"
)

type LessonProgress struct {
	ID              int       `json:"id"`
	UserID          int       `json:"user_id"`
	LessonID        int       `json:"lesson_id"`
	CourseID        int       `json:"course_id"`
    IsCompleted     bool      `json:"is_completed"`    // Прочитано или нет
	CompletedAt     time.Time `json:"completed_at"`    // Время прочтения
	CreatedAt       time.Time `json:"created_at"` // Время создания записи
	UpdatedAt       time.Time `json:"updated_at"` // Время обновления записи

}