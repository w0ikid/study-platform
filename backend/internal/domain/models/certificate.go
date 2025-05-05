package models

import "time"

type Certificate struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	CourseID  int       `json:"course_id"`
	IssuedAt  time.Time `json:"issued_at"`
}