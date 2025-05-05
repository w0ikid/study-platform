package models

import "time"
type Enrollment struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	CourseID  int       `json:"course_id"`
	Status    string    `json:"status"` // roles: active, completed
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}