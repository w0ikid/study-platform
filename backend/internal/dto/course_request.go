package dto

type CreateCourseRequest struct {
	Name string `json:"name" validate:"required"`
	Description string `json:"description" validate:"required"`
	ImageUrl string `json:"image_url"`
}
