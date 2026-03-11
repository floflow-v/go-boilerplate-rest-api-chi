package author_test

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"go-boilerplate-rest-api-chi/internal/author"
	"go-boilerplate-rest-api-chi/internal/author/dto"
	"go-boilerplate-rest-api-chi/internal/database/sqlc"
	internalError "go-boilerplate-rest-api-chi/internal/error"
	"go-boilerplate-rest-api-chi/internal/mocks"
)

func TestAuthorService_CreateAuthor(t *testing.T) {
	tests := []struct {
		name             string
		input            dto.CreateAuthorRequest
		configureMock    func(*mocks.MockQuerier)
		expectedResponse dto.AuthorResponse
		expectedError    error
	}{
		{
			name: "success create author",
			input: dto.CreateAuthorRequest{
				Name: "J.K. Rowling",
			},
			configureMock: func(mockQuerier *mocks.MockQuerier) {
				mockQuerier.EXPECT().
					CreateAuthor(gomock.Any(), gomock.AssignableToTypeOf(sqlc.CreateAuthorParams{})).
					DoAndReturn(func(ctx context.Context, params sqlc.CreateAuthorParams) error {
						assert.Equal(t, "J.K. Rowling", params.Name)
						mockQuerier.EXPECT().
							GetAuthorByID(gomock.Any(), params.ID).
							Return(sqlc.Author(params), nil)
						return nil
					})

			},
			expectedResponse: dto.AuthorResponse{
				Name: "J.K. Rowling",
			},
		},
		{
			name: "error duplicate author",
			input: dto.CreateAuthorRequest{
				Name: "Duplicate Author",
			},
			configureMock: func(mockQuerier *mocks.MockQuerier) {
				mockQuerier.EXPECT().
					CreateAuthor(gomock.Any(), gomock.AssignableToTypeOf(sqlc.CreateAuthorParams{})).
					Return(&mysql.MySQLError{Number: 1062})
			},
			expectedError: internalError.AuthorDuplicate,
		},
		{
			name: "error database error",
			input: dto.CreateAuthorRequest{
				Name: "Test Author",
			},
			configureMock: func(mockQuerier *mocks.MockQuerier) {
				mockQuerier.EXPECT().
					CreateAuthor(gomock.Any(), gomock.AssignableToTypeOf(sqlc.CreateAuthorParams{})).
					Return(errors.New("invalid db"))
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
			service := author.NewAuthorService(mockQuerier, zerolog.Nop())

			result, err := service.CreateAuthor(context.Background(), test.input)

			if test.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, test.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, test.expectedResponse.Name, result.Name)
				assert.NotEqual(t, uuid.Nil, result.ID)
			}
		})
	}
}

func TestAuthorService_GetAuthorByID(t *testing.T) {
	tests := []struct {
		name             string
		authorID         uuid.UUID
		configureMock    func(*mocks.MockQuerier)
		expectedResponse dto.AuthorResponse
		expectedError    error
	}{
		{
			name:     "success get author by id",
			authorID: uuid.MustParse("eb21d07a-7ab3-40db-bfd3-448093bc5626"),
			configureMock: func(mockQuerier *mocks.MockQuerier) {
				mockQuerier.EXPECT().
					GetAuthorByID(gomock.Any(), "eb21d07a-7ab3-40db-bfd3-448093bc5626").
					Return(sqlc.Author{
						ID:   "eb21d07a-7ab3-40db-bfd3-448093bc5626",
						Name: "J.K. Rowling",
					}, nil)
			},
			expectedResponse: dto.AuthorResponse{
				ID:   "eb21d07a-7ab3-40db-bfd3-448093bc5626",
				Name: "J.K. Rowling",
			},
		},
		{
			name:     "error author not found",
			authorID: uuid.MustParse("eb21d07a-7ab3-40db-bfd3-448093bc5626"),
			configureMock: func(mockQuerier *mocks.MockQuerier) {
				mockQuerier.EXPECT().
					GetAuthorByID(gomock.Any(), "eb21d07a-7ab3-40db-bfd3-448093bc5626").
					Return(sqlc.Author{}, sql.ErrNoRows)
			},
			expectedError: internalError.AuthorNotFound,
		},
		{
			name:     "error database error",
			authorID: uuid.MustParse("eb21d07a-7ab3-40db-bfd3-448093bc5626"),
			configureMock: func(mockQuerier *mocks.MockQuerier) {
				mockQuerier.EXPECT().
					GetAuthorByID(gomock.Any(), "eb21d07a-7ab3-40db-bfd3-448093bc5626").
					Return(sqlc.Author{}, errors.New("invalid db"))
			},
			expectedError: internalError.InternalError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			t.Cleanup(ctrl.Finish)
			authorRepoMock := mocks.NewMockQuerier(ctrl)

			test.configureMock(authorRepoMock)
			service := author.NewAuthorService(authorRepoMock, zerolog.Nop())

			result, err := service.GetAuthorByID(context.Background(), test.authorID)

			if test.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, test.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, test.expectedResponse.Name, result.Name)
				assert.NotEqual(t, uuid.Nil, result.ID)
			}
		})
	}
}
