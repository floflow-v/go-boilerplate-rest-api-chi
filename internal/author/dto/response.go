package dto

import (
	"go-rest-api-chi-example/internal/database/sqlc"
)

type AuthorResponse struct {
	ID   string `json:"id" example:"019cbe7c-01f2-7b7e-8424-818397b8652c"`
	Name string `json:"name" example:"J.K. Rowling"`
}

func ToAuthorResponse(author sqlc.Author) AuthorResponse {
	return AuthorResponse{
		ID:   author.ID,
		Name: author.Name,
	}
}
