package dto

import (
	authorDTO "go-boilerplate-rest-api-chi/internal/author/dto"
	"go-boilerplate-rest-api-chi/internal/model"
)

type BookResponse struct {
	ID          string                   `json:"id" example:"019cbe7c-01f2-7b7e-8424-818397b8652c"`
	Title       string                   `json:"title" example:"Harry Potter and the Philosopher's Stone"`
	Description string                   `json:"description" example:"Harry Potter has never even heard of Hogwarts when the letters start dropping on the doormat at number four, Privet Drive. Addressed in green ink on yellowish parchment with a purple seal, they are swiftly confiscated by his grisly aunt and uncle."`
	Author      authorDTO.AuthorResponse `json:"author"`
}

func ToBookResponse(book *model.Book) *BookResponse {
	response := &BookResponse{
		ID:          book.ID.String(),
		Title:       book.Title,
		Description: book.Description,
		Author: authorDTO.AuthorResponse{
			ID:   book.Author.ID.String(),
			Name: book.Author.Name,
		},
	}

	return response
}

func ToBooksResponse(books []*model.Book) []BookResponse {
	responses := make([]BookResponse, len(books))
	for i, book := range books {
		responses[i] = *ToBookResponse(book)
	}
	return responses
}
