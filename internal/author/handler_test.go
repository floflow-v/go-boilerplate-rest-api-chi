package author_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"go-rest-api-chi-example/internal/author"
	"go-rest-api-chi-example/internal/author/dto"
	internalError "go-rest-api-chi-example/internal/error"
	"go-rest-api-chi-example/internal/mocks"
	"go-rest-api-chi-example/internal/validator"
)

func TestAuthorHandler_CreateAuthor(t *testing.T) {
	tests := []struct {
		name               string
		requestBody        interface{}
		configureMock      func(*mocks.MockAuthorService)
		expectedStatusCode int
		expectedResponse   interface{}
	}{
		{
			name: "success create author",
			requestBody: dto.CreateAuthorRequest{
				Name: "George R.R. Martin",
			},
			configureMock: func(mockService *mocks.MockAuthorService) {
				input := &dto.CreateAuthorRequest{
					Name: "George R.R. Martin",
				}

				mockService.EXPECT().
					CreateAuthor(gomock.Any(), input).
					Return(&dto.AuthorResponse{
						ID:   "aeca0955-bae4-47e9-9f85-6818dc68ca51",
						Name: "George R.R. Martin",
					}, nil)
			},
			expectedStatusCode: http.StatusCreated,
			expectedResponse: dto.AuthorResponse{
				ID:   "aeca0955-bae4-47e9-9f85-6818dc68ca51",
				Name: "George R.R. Martin",
			},
		},
		{
			name:               "error invalid JSON",
			requestBody:        nil,
			configureMock:      func(mockService *mocks.MockAuthorService) {},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse: internalError.Response{
				Error:   internalError.InvalidRequestBody.Code,
				Message: internalError.InvalidRequestBody.Message,
			},
		},
		{
			name: "error validation fails empty name",
			requestBody: dto.CreateAuthorRequest{
				Name: "",
			},
			configureMock:      func(mockService *mocks.MockAuthorService) {},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse: internalError.Response{
				Error:   internalError.ValidationError.Code,
				Message: internalError.ValidationError.Message,
				Details: []interface{}{
					map[string]interface{}{
						"field":   "name",
						"message": "name is required",
					},
				},
			},
		},
		{
			name: "error duplicate author",
			requestBody: dto.CreateAuthorRequest{
				Name: "Duplicate Author",
			},
			configureMock: func(mockService *mocks.MockAuthorService) {
				input := &dto.CreateAuthorRequest{
					Name: "Duplicate Author",
				}

				mockService.EXPECT().
					CreateAuthor(gomock.Any(), input).
					Return(nil, internalError.AuthorDuplicate)
			},
			expectedStatusCode: http.StatusConflict,
			expectedResponse: internalError.Response{
				Error:   internalError.AuthorDuplicate.Code,
				Message: internalError.AuthorDuplicate.Message,
			},
		},
		{
			name: "error service internal error",
			requestBody: dto.CreateAuthorRequest{
				Name: "George R.R. Martin",
			},
			configureMock: func(mockService *mocks.MockAuthorService) {
				input := &dto.CreateAuthorRequest{
					Name: "George R.R. Martin",
				}

				mockService.EXPECT().
					CreateAuthor(gomock.Any(), input).
					Return(nil, errors.New("database connection failed"))
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedResponse: internalError.Response{
				Error:   internalError.InternalError.Code,
				Message: internalError.InternalError.Message,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			t.Cleanup(ctrl.Finish)

			mockService := mocks.NewMockAuthorService(ctrl)
			test.configureMock(mockService)

			v := validator.New()
			handler := author.NewAuthorHandler(mockService, v, zerolog.Nop())

			var body *bytes.Buffer
			if test.requestBody == nil {
				body = bytes.NewBuffer([]byte{})
			} else {
				b, err := json.Marshal(test.requestBody)
				require.NoError(t, err)
				body = bytes.NewBuffer(b)
			}

			req := httptest.NewRequest(http.MethodPost, "/authors", body)
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			r := chi.NewRouter()
			r.Mount("/authors", handler.Routes())

			r.ServeHTTP(w, req)

			assert.Equal(t, test.expectedStatusCode, w.Code)

			expectedJSON, err := json.Marshal(test.expectedResponse)
			require.NoError(t, err)

			assert.JSONEq(t, string(expectedJSON), w.Body.String())

		})
	}
}

func TestAuthorHandler_GetAuthorByID(t *testing.T) {
	tests := []struct {
		name               string
		idInUrlParam       string
		configureMock      func(*mocks.MockAuthorService)
		expectedStatusCode int
		expectedResponse   interface{}
	}{
		{
			name:         "success get author by ID",
			idInUrlParam: "aeca0955-bae4-47e9-9f85-6818dc68ca51",
			configureMock: func(mockService *mocks.MockAuthorService) {
				mockService.EXPECT().
					GetAuthorByID(gomock.Any(), uuid.MustParse("aeca0955-bae4-47e9-9f85-6818dc68ca51")).
					Return(&dto.AuthorResponse{
						ID:   "aeca0955-bae4-47e9-9f85-6818dc68ca51",
						Name: "George R.R. Martin",
					}, nil)
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse: dto.AuthorResponse{
				ID:   "aeca0955-bae4-47e9-9f85-6818dc68ca51",
				Name: "George R.R. Martin",
			},
		},
		{
			name:               "error invalid uuid",
			idInUrlParam:       "invalid-uuid",
			configureMock:      func(mockService *mocks.MockAuthorService) {},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse: internalError.Response{
				Error:   internalError.InvalidUUID.Code,
				Message: internalError.InvalidUUID.Message,
			},
		},
		{
			name:         "error author not found",
			idInUrlParam: "aeca0955-bae4-47e9-9f85-6818dc68ca51",
			configureMock: func(mockService *mocks.MockAuthorService) {
				mockService.EXPECT().
					GetAuthorByID(gomock.Any(), uuid.MustParse("aeca0955-bae4-47e9-9f85-6818dc68ca51")).
					Return(nil, internalError.AuthorNotFound)
			},
			expectedStatusCode: http.StatusNotFound,
			expectedResponse: internalError.Response{
				Error:   internalError.AuthorNotFound.Code,
				Message: internalError.AuthorNotFound.Message,
			},
		},
		{
			name:         "error service internal error",
			idInUrlParam: "aeca0955-bae4-47e9-9f85-6818dc68ca51",
			configureMock: func(mockService *mocks.MockAuthorService) {
				mockService.EXPECT().
					GetAuthorByID(gomock.Any(), uuid.MustParse("aeca0955-bae4-47e9-9f85-6818dc68ca51")).
					Return(nil, errors.New("database connection failed"))
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedResponse: internalError.Response{
				Error:   internalError.InternalError.Code,
				Message: internalError.InternalError.Message,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			t.Cleanup(ctrl.Finish)

			mockService := mocks.NewMockAuthorService(ctrl)
			test.configureMock(mockService)

			v := validator.New()
			handler := author.NewAuthorHandler(mockService, v, zerolog.Nop())

			url := fmt.Sprintf("/authors/%s", test.idInUrlParam)
			req := httptest.NewRequest(http.MethodGet, url, nil)
			w := httptest.NewRecorder()

			r := chi.NewRouter()
			r.Mount("/authors", handler.Routes())

			r.ServeHTTP(w, req)

			assert.Equal(t, test.expectedStatusCode, w.Code)

			expectedJSON, err := json.Marshal(test.expectedResponse)
			require.NoError(t, err)

			assert.JSONEq(t, string(expectedJSON), w.Body.String())
		})
	}
}
