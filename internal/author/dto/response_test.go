package dto_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"go-rest-api-chi-example/internal/author/dto"
	"go-rest-api-chi-example/internal/model"
)

func TestToAuthorResponse(t *testing.T) {
	t.Run("nominal", func(t *testing.T) {
		entity := model.Author{
			ID:   uuid.MustParse("aeca0955-bae4-47e9-9f85-6818dc68ca51"),
			Name: "George R.R. Martin",
		}

		expectedResponse := dto.AuthorResponse{
			ID:   "aeca0955-bae4-47e9-9f85-6818dc68ca51",
			Name: "George R.R. Martin",
		}

		response := dto.ToAuthorResponse(&entity)

		assert.Equal(t, &expectedResponse, response)
	})
}
