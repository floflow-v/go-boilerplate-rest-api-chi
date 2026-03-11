package dto

import (
	authorDTO "go-boilerplate-rest-api-chi/internal/author/dto"
	"go-boilerplate-rest-api-chi/internal/database/sqlc"
)

type BookResponse struct {
	ID          string                   `json:"id" example:"019cbe7c-01f2-7b7e-8424-818397b8652c"`
	Title       string                   `json:"title" example:"Harry Potter and the Philosopher's Stone"`
	Description string                   `json:"description" example:"Harry Potter has never even heard of Hogwarts when the letters start dropping on the doormat at number four, Privet Drive. Addressed in green ink on yellowish parchment with a purple seal, they are swiftly confiscated by his grisly aunt and uncle."`
	Author      authorDTO.AuthorResponse `json:"author"`
}

func ToBookResponse(book sqlc.GetBookByIDRow) BookResponse {
	return BookResponse{
		ID:          book.ID,
		Title:       book.Title,
		Description: book.Description,
		Author: authorDTO.AuthorResponse{
			ID:   book.AuthorID,
			Name: book.AuthorName,
		},
	}
}

func ToBookResponseFromRows(books []sqlc.GetAllBooksRow) []BookResponse {
	responses := make([]BookResponse, len(books))

	for i, book := range books {
		responses[i] = BookResponse{
			ID:          book.ID,
			Title:       book.Title,
			Description: book.Description,
			Author: authorDTO.AuthorResponse{
				ID:   book.AuthorID,
				Name: book.AuthorName,
			},
		}
	}
	return responses
}
