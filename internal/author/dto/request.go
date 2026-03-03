package dto

type CreateAuthorRequest struct {
	Name string `json:"name" example:"J.K. Rowling" validate:"required"`
}
