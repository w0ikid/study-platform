package models

import (
	"time"
)

type Lesson struct {
    ID        int       `json:"id"`
    CourseID  int       `json:"course_id"`
    Title     string    `json:"title"`
    Content   string    `json:"content"`
    VideoURL  string    `json:"video_url,omitempty"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}
// 