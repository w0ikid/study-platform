package dto

type CreateLessonRequest struct {
	Title    string `json:"title" validate:"required"`
	Content  string `json:"content" validate:"required"`
	VideoURL string `json:"video_url,omitempty" validate:"omitempty,url"`
}