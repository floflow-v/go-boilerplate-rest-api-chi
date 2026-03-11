package dto_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	authorDTO "go-rest-api-chi-example/internal/author/dto"
	"go-rest-api-chi-example/internal/book/dto"
	"go-rest-api-chi-example/internal/database/sqlc"
)

func ToBookResponse(t *testing.T) {
	t.Run("nominal", func(t *testing.T) {
		entity := sqlc.GetBookByIDRow{
			ID:          uuid.MustParse("58411bf8-aa11-4553-9b13-4bdf58875d35").String(),
			Title:       "Book1",
			Description: "Description1",
			AuthorID:    uuid.MustParse("b846fc59-401a-450d-b3f1-3e9a953d7c22").String(),
			AuthorName:  "Author1",
		}

		expectedResponse := dto.BookResponse{
			ID:          "58411bf8-aa11-4553-9b13-4bdf58875d35",
			Title:       "Book1",
			Description: "Description1",
			Author: authorDTO.AuthorResponse{
				ID:   "b846fc59-401a-450d-b3f1-3e9a953d7c22",
				Name: "Author1",
			},
		}

		response := dto.ToBookResponse(entity)

		assert.Equal(t, &expectedResponse, response)
	})
}

func ToBookResponseFromRows(t *testing.T) {
	t.Run("nominal for multiple books", func(t *testing.T) {
		entities := []sqlc.GetAllBooksRow{
			{
				ID:          uuid.MustParse("58411bf8-aa11-4553-9b13-4bdf58875d35").String(),
				Title:       "Book1",
				Description: "Description1",
				AuthorID:    uuid.MustParse("b846fc59-401a-450d-b3f1-3e9a953d7c22").String(),
				AuthorName:  "Author1",
			},
			{
				ID:          uuid.MustParse("26dcbb76-e09d-45c5-8c97-f32ce3a2766a").String(),
				Title:       "Book2",
				Description: "Description2",
				AuthorID:    uuid.MustParse("b846fc59-401a-450d-b3f1-3e9a953d7c22").String(),
				AuthorName:  "Author1",
			},
			{
				ID:          uuid.MustParse("2933e943-bde9-4743-8961-97828d166e11").String(),
				Title:       "Book3",
				Description: "Description3",
				AuthorID:    uuid.MustParse("b846fc59-401a-450d-b3f1-3e9a953d7c22").String(),
				AuthorName:  "Author1",
			},
		}

		expectedResponse := []dto.BookResponse{
			{
				ID:          "58411bf8-aa11-4553-9b13-4bdf58875d35",
				Title:       "Book1",
				Description: "Description1",
				Author: authorDTO.AuthorResponse{
					ID:   "b846fc59-401a-450d-b3f1-3e9a953d7c22",
					Name: "Author1",
				},
			},
			{
				ID:          "26dcbb76-e09d-45c5-8c97-f32ce3a2766a",
				Title:       "Book2",
				Description: "Description2",
				Author: authorDTO.AuthorResponse{
					ID:   "b846fc59-401a-450d-b3f1-3e9a953d7c22",
					Name: "Author1",
				},
			},
			{
				ID:          "2933e943-bde9-4743-8961-97828d166e11",
				Title:       "Book3",
				Description: "Description3",
				Author: authorDTO.AuthorResponse{
					ID:   "b846fc59-401a-450d-b3f1-3e9a953d7c22",
					Name: "Author1",
				},
			},
		}

		response := dto.ToBookResponseFromRows(entities)

		assert.Equal(t, expectedResponse, response)
	})
}
