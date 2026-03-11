package book_test

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	authorDto "go-boilerplate-rest-api-chi/internal/author/dto"
	"go-boilerplate-rest-api-chi/internal/book"
	"go-boilerplate-rest-api-chi/internal/book/dto"
	"go-boilerplate-rest-api-chi/internal/database/sqlc"
	internalError "go-boilerplate-rest-api-chi/internal/error"
	"go-boilerplate-rest-api-chi/internal/mocks"
)

func TestBookService_CreateBook(t *testing.T) {
	tests := []struct {
		name             string
		input            dto.CreateBookRequest
		configureMock    func(*mocks.MockQuerier)
		expectedResponse dto.BookResponse
		expectedError    error
	}{
		{
			name: "success create book",
			input: dto.CreateBookRequest{
				Title:       "Book1",
				Description: "Description of book1",
				AuthorID:    "779404e4-2660-4c80-b958-cfa72515e7d4",
			},
			configureMock: func(mockQuerier *mocks.MockQuerier) {
				authorID := uuid.MustParse("779404e4-2660-4c80-b958-cfa72515e7d4")

				mockQuerier.EXPECT().
					GetAuthorByID(gomock.Any(), authorID.String()).
					Return(sqlc.Author{
						ID:   authorID.String(),
						Name: "Author1",
					}, nil)

				mockQuerier.EXPECT().
					CreateBook(gomock.Any(), gomock.Any()).
					Return(nil)

				mockQuerier.EXPECT().
					GetBookByID(gomock.Any(), gomock.Any()).
					Return(sqlc.GetBookByIDRow{
						ID:          uuid.New().String(),
						Title:       "Book1",
						Description: "Description of book1",
						AuthorID:    authorID.String(),
						AuthorName:  "Author1",
					}, nil)
			},
			expectedResponse: dto.BookResponse{
				Title:       "Book1",
				Description: "Description of book1",
				Author: authorDto.AuthorResponse{
					ID:   "779404e4-2660-4c80-b958-cfa72515e7d4",
					Name: "Author1",
				},
			},
		},
		{
			name: "error invalid AuthorID",
			input: dto.CreateBookRequest{
				Title:       "Book1",
				Description: "Description of book1",
				AuthorID:    "invalid-uuid",
			},
			configureMock: func(mockQuerier *mocks.MockQuerier) {},
			expectedError: internalError.InvalidAuthorID,
		},
		{
			name: "error author not found",
			input: dto.CreateBookRequest{
				Title:       "Book1",
				Description: "Description of book1",
				AuthorID:    "779404e4-2660-4c80-b958-cfa72515e7d4",
			},
			configureMock: func(mockQuerier *mocks.MockQuerier) {
				authorID := uuid.MustParse("779404e4-2660-4c80-b958-cfa72515e7d4")

				mockQuerier.EXPECT().
					GetAuthorByID(gomock.Any(), authorID.String()).
					Return(sqlc.Author{}, sql.ErrNoRows)
			},
			expectedError: internalError.AuthorNotFound,
		},
		{
			name: "error database error on GetAuthorByID",
			input: dto.CreateBookRequest{
				Title:       "Book1",
				Description: "Description of book1",
				AuthorID:    "779404e4-2660-4c80-b958-cfa72515e7d4",
			},
			configureMock: func(mockQuerier *mocks.MockQuerier) {
				authorID := uuid.MustParse("779404e4-2660-4c80-b958-cfa72515e7d4")

				mockQuerier.EXPECT().
					GetAuthorByID(gomock.Any(), authorID.String()).
					Return(sqlc.Author{}, errors.New("db connection failed"))
			},
			expectedError: internalError.InternalError,
		},
		{
			name: "error database error on CreateBook",
			input: dto.CreateBookRequest{
				Title:       "Book1",
				Description: "Description of book1",
				AuthorID:    "779404e4-2660-4c80-b958-cfa72515e7d4",
			},
			configureMock: func(mockQuerier *mocks.MockQuerier) {
				authorID := uuid.MustParse("779404e4-2660-4c80-b958-cfa72515e7d4")

				mockQuerier.EXPECT().
					GetAuthorByID(gomock.Any(), authorID.String()).
					Return(sqlc.Author{
						ID:   authorID.String(),
						Name: "Author1",
					}, nil)

				mockQuerier.EXPECT().
					CreateBook(gomock.Any(), gomock.Any()).
					Return(errors.New("db connection failed"))
			},
			expectedError: internalError.InternalError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			t.Cleanup(ctrl.Finish)

			mockQuerier := mocks.NewMockQuerier(ctrl)
			test.configureMock(mockQuerier)

			service := book.NewBookService(mockQuerier, zerolog.Nop())

			result, err := service.CreateBook(context.Background(), test.input)

			if test.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, test.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, result.ID)
				assert.Equal(t, test.expectedResponse.Title, result.Title)
				assert.Equal(t, test.expectedResponse.Description, result.Description)
				assert.Equal(t, test.expectedResponse.Author.ID, result.Author.ID)
				assert.Equal(t, test.expectedResponse.Author.Name, result.Author.Name)
			}
		})
	}
}

func TestBookService_GetAllBooks(t *testing.T) {
	tests := []struct {
		name             string
		configureMock    func(*mocks.MockQuerier)
		expectedResponse []dto.BookResponse
		expectedError    error
	}{
		{
			name: "success get all books",
			configureMock: func(mockQuerier *mocks.MockQuerier) {
				authorID := uuid.MustParse("779404e4-2660-4c80-b958-cfa72515e7d4")

				mockQuerier.EXPECT().
					GetAllBooks(gomock.Any()).
					Return([]sqlc.GetAllBooksRow{
						{
							ID:          "619c69fb-9bcb-451e-b825-29b81697a531",
							Title:       "Book1",
							Description: "Description of book1",
							AuthorID:    authorID.String(),
							AuthorName:  "Author1",
						},
						{
							ID:          "92e37eac-0097-405e-a5cc-7055a39bdde8",
							Title:       "Book2",
							Description: "Description of book2",
							AuthorID:    authorID.String(),
							AuthorName:  "Author1",
						},
						{
							ID:          "dec6b550-292e-45ad-b269-f59bfa06bf01",
							Title:       "Book3",
							Description: "Description of book3",
							AuthorID:    authorID.String(),
							AuthorName:  "Author1",
						},
					}, nil)
			},
			expectedResponse: []dto.BookResponse{
				{
					ID:          "619c69fb-9bcb-451e-b825-29b81697a531",
					Title:       "Book1",
					Description: "Description of book1",
					Author:      authorDto.AuthorResponse{ID: "779404e4-2660-4c80-b958-cfa72515e7d4", Name: "Author1"},
				},
				{
					ID:          "92e37eac-0097-405e-a5cc-7055a39bdde8",
					Title:       "Book2",
					Description: "Description of book2",
					Author:      authorDto.AuthorResponse{ID: "779404e4-2660-4c80-b958-cfa72515e7d4", Name: "Author1"},
				},
				{
					ID:          "dec6b550-292e-45ad-b269-f59bfa06bf01",
					Title:       "Book3",
					Description: "Description of book3",
					Author:      authorDto.AuthorResponse{ID: "779404e4-2660-4c80-b958-cfa72515e7d4", Name: "Author1"},
				},
			},
		},
		{
			name: "error no books found",
			configureMock: func(mockQuerier *mocks.MockQuerier) {
				mockQuerier.EXPECT().
					GetAllBooks(gomock.Any()).
					Return([]sqlc.GetAllBooksRow{}, nil)
			},
			expectedError: internalError.BookNotFound,
		},
		{
			name: "error database connection failed",
			configureMock: func(mockQuerier *mocks.MockQuerier) {
				mockQuerier.EXPECT().
					GetAllBooks(gomock.Any()).
					Return(nil, errors.New("db connection failed"))
			},
			expectedError: internalError.InternalError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			t.Cleanup(ctrl.Finish)

			mockQuerier := mocks.NewMockQuerier(ctrl)
			test.configureMock(mockQuerier)

			service := book.NewBookService(mockQuerier, zerolog.Nop())

			books, err := service.GetAllBooks(context.Background())

			if test.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, test.expectedError.Error(), err.Error())
				assert.Nil(t, books)
			} else {
				assert.NoError(t, err)
				assert.Len(t, books, len(test.expectedResponse))
				for i := range books {
					assert.Equal(t, test.expectedResponse[i].ID, books[i].ID)
					assert.Equal(t, test.expectedResponse[i].Title, books[i].Title)
					assert.Equal(t, test.expectedResponse[i].Description, books[i].Description)
					assert.Equal(t, test.expectedResponse[i].Author.ID, books[i].Author.ID)
					assert.Equal(t, test.expectedResponse[i].Author.Name, books[i].Author.Name)
				}
			}
		})
	}
}

func TestBookService_GetBookByID(t *testing.T) {
	tests := []struct {
		name             string
		bookID           uuid.UUID
		configureMock    func(*mocks.MockQuerier)
		expectedResponse dto.BookResponse
		expectedError    error
	}{
		{
			name:   "success get book by ID",
			bookID: uuid.MustParse("c6efb683-455d-4ca0-b8aa-83ca9b930a06"),
			configureMock: func(mockQuerier *mocks.MockQuerier) {
				mockQuerier.EXPECT().
					GetBookByID(gomock.Any(), "c6efb683-455d-4ca0-b8aa-83ca9b930a06").
					Return(sqlc.GetBookByIDRow{
						ID:          "c6efb683-455d-4ca0-b8aa-83ca9b930a06",
						Title:       "Book1",
						Description: "Description of book1",
						AuthorID:    "b57063a5-409c-457c-bbec-d5850b2e3761",
						AuthorName:  "Author1",
					}, nil)
			},
			expectedResponse: dto.BookResponse{
				ID:          "c6efb683-455d-4ca0-b8aa-83ca9b930a06",
				Title:       "Book1",
				Description: "Description of book1",
				Author:      authorDto.AuthorResponse{ID: "b57063a5-409c-457c-bbec-d5850b2e3761", Name: "Author1"},
			},
		},
		{
			name:   "error book not found",
			bookID: uuid.MustParse("c6efb683-455d-4ca0-b8aa-83ca9b930a06"),
			configureMock: func(mockQuerier *mocks.MockQuerier) {
				mockQuerier.EXPECT().
					GetBookByID(gomock.Any(), "c6efb683-455d-4ca0-b8aa-83ca9b930a06").
					Return(sqlc.GetBookByIDRow{}, sql.ErrNoRows)
			},
			expectedError: internalError.BookNotFound,
		},
		{
			name:   "error database connection failed",
			bookID: uuid.MustParse("c6efb683-455d-4ca0-b8aa-83ca9b930a06"),
			configureMock: func(mockQuerier *mocks.MockQuerier) {
				mockQuerier.EXPECT().
					GetBookByID(gomock.Any(), "c6efb683-455d-4ca0-b8aa-83ca9b930a06").
					Return(sqlc.GetBookByIDRow{}, errors.New("db connection failed"))
			},
			expectedError: internalError.InternalError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			t.Cleanup(ctrl.Finish)

			mockQuerier := mocks.NewMockQuerier(ctrl)
			test.configureMock(mockQuerier)

			service := book.NewBookService(mockQuerier, zerolog.Nop())

			result, err := service.GetBookByID(context.Background(), test.bookID)

			if test.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, test.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.expectedResponse.ID, result.ID)
				assert.Equal(t, test.expectedResponse.Title, result.Title)
				assert.Equal(t, test.expectedResponse.Description, result.Description)
				assert.Equal(t, test.expectedResponse.Author.ID, result.Author.ID)
				assert.Equal(t, test.expectedResponse.Author.Name, result.Author.Name)
			}
		})
	}
}

func TestBookService_UpdateBook(t *testing.T) {
	tests := []struct {
		name          string
		bookID        uuid.UUID
		input         dto.UpdateBookRequest
		configureMock func(*mocks.MockQuerier)
		expectedError error
	}{
		{
			name:   "success update book",
			bookID: uuid.MustParse("c6efb683-455d-4ca0-b8aa-83ca9b930a06"),
			input: dto.UpdateBookRequest{
				Title:       "Updated Title",
				Description: "Updated description",
			},
			configureMock: func(mockQuerier *mocks.MockQuerier) {
				mockQuerier.EXPECT().
					GetBookByID(gomock.Any(), "c6efb683-455d-4ca0-b8aa-83ca9b930a06").
					Return(sqlc.GetBookByIDRow{
						ID:          "c6efb683-455d-4ca0-b8aa-83ca9b930a06",
						Title:       "Title",
						Description: "description",
					}, nil)

				mockQuerier.EXPECT().
					UpdateBook(gomock.Any(), sqlc.UpdateBookParams{
						Title:       "Updated Title",
						Description: "Updated description",
						ID:          "c6efb683-455d-4ca0-b8aa-83ca9b930a06",
					}).
					Return(nil)

			},
		},
		{
			name:   "error book not found",
			bookID: uuid.MustParse("c6efb683-455d-4ca0-b8aa-83ca9b930a06"),
			input: dto.UpdateBookRequest{
				Title:       "Updated Title",
				Description: "Updated description",
			},
			configureMock: func(mockQuerier *mocks.MockQuerier) {
				mockQuerier.EXPECT().
					GetBookByID(gomock.Any(), "c6efb683-455d-4ca0-b8aa-83ca9b930a06").
					Return(sqlc.GetBookByIDRow{}, sql.ErrNoRows)
			},
			expectedError: internalError.BookNotFound,
		},
		{
			name:   "error database connection failed",
			bookID: uuid.MustParse("c6efb683-455d-4ca0-b8aa-83ca9b930a06"),
			input: dto.UpdateBookRequest{
				Title:       "Updated Title",
				Description: "Updated description",
			},
			configureMock: func(mockQuerier *mocks.MockQuerier) {
				mockQuerier.EXPECT().
					GetBookByID(gomock.Any(), gomock.Any()).
					Return(sqlc.GetBookByIDRow{}, errors.New("db connection failed"))
			},
			expectedError: internalError.InternalError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			t.Cleanup(ctrl.Finish)

			mockQuerier := mocks.NewMockQuerier(ctrl)
			test.configureMock(mockQuerier)

			service := book.NewBookService(mockQuerier, zerolog.Nop())

			err := service.UpdateBook(context.Background(), test.input, test.bookID)

			if test.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, test.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestBookService_DeleteBook(t *testing.T) {
	tests := []struct {
		name          string
		bookID        uuid.UUID
		configureMock func(*mocks.MockQuerier)
		expectedError error
	}{
		{
			name:   "success delete book",
			bookID: uuid.MustParse("c6efb683-455d-4ca0-b8aa-83ca9b930a06"),
			configureMock: func(mockQuerier *mocks.MockQuerier) {
				mockQuerier.EXPECT().
					GetBookByID(gomock.Any(), "c6efb683-455d-4ca0-b8aa-83ca9b930a06").
					Return(sqlc.GetBookByIDRow{
						ID:          "c6efb683-455d-4ca0-b8aa-83ca9b930a06",
						Title:       "Title",
						Description: "description",
					}, nil)

				mockQuerier.EXPECT().
					DeleteBook(gomock.Any(), "c6efb683-455d-4ca0-b8aa-83ca9b930a06").
					Return(nil)
			},
		},
		{
			name:   "error book not found",
			bookID: uuid.MustParse("c6efb683-455d-4ca0-b8aa-83ca9b930a06"),
			configureMock: func(mockQuerier *mocks.MockQuerier) {
				mockQuerier.EXPECT().
					GetBookByID(gomock.Any(), "c6efb683-455d-4ca0-b8aa-83ca9b930a06").
					Return(sqlc.GetBookByIDRow{}, sql.ErrNoRows)
			},
			expectedError: internalError.BookNotFound,
		},
		{
			name:   "error database connection failed",
			bookID: uuid.MustParse("c6efb683-455d-4ca0-b8aa-83ca9b930a06"),
			configureMock: func(mockQuerier *mocks.MockQuerier) {
				mockQuerier.EXPECT().
					GetBookByID(gomock.Any(), "c6efb683-455d-4ca0-b8aa-83ca9b930a06").
					Return(sqlc.GetBookByIDRow{}, errors.New("db connection failed"))
			},
			expectedError: internalError.InternalError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			t.Cleanup(ctrl.Finish)

			mockQuerier := mocks.NewMockQuerier(ctrl)
			test.configureMock(mockQuerier)

			service := book.NewBookService(mockQuerier, zerolog.Nop())

			err := service.DeleteBook(context.Background(), test.bookID)

			if test.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, test.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
