package author

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/rs/zerolog"

	"go-rest-api-chi-example/internal/author/dto"
	"go-rest-api-chi-example/internal/database/sqlc"
	internalError "go-rest-api-chi-example/internal/error"
)

//go:generate mockgen -destination=../mocks/mock_author_service.go -package=mocks go-rest-api-chi-example/internal/author AuthorService
type AuthorService interface {
	CreateAuthor(ctx context.Context, req dto.CreateAuthorRequest) (dto.AuthorResponse, error)
	GetAuthorByID(ctx context.Context, authorID uuid.UUID) (dto.AuthorResponse, error)
}

type authorService struct {
	querier sqlc.Querier
	logger  zerolog.Logger
}

func NewAuthorService(querier sqlc.Querier, logger zerolog.Logger) AuthorService {
	return &authorService{
		querier: querier,
		logger:  logger,
	}
}

func (s *authorService) CreateAuthor(ctx context.Context, req dto.CreateAuthorRequest) (dto.AuthorResponse, error) {
	authorID, _ := uuid.NewV7()

	authorRequest := sqlc.CreateAuthorParams{
		ID:   authorID.String(),
		Name: req.Name,
	}

	err := s.querier.CreateAuthor(ctx, authorRequest)
	if err != nil {
		s.logger.Error().Err(err).Msg("Failed to create author")

		dbErr := internalError.MapDBError(err)

		switch {
		case errors.Is(dbErr, internalError.ErrDBErrDuplicate):
			return dto.AuthorResponse{}, internalError.AuthorDuplicate

		default:
			s.logger.Error().Err(err).Msg("Unknown error")
			return dto.AuthorResponse{}, internalError.InternalError
		}
	}

	author, err := s.querier.GetAuthorByID(ctx, authorID.String())
	if err != nil {
		s.logger.Error().Err(err).Msg("failed to get author after creation")
		return dto.AuthorResponse{}, internalError.InternalError
	}

	return dto.ToAuthorResponse(author), nil
}

func (s *authorService) GetAuthorByID(ctx context.Context, authorID uuid.UUID) (dto.AuthorResponse, error) {
	author, err := s.querier.GetAuthorByID(ctx, authorID.String())
	if err != nil {
		dbErr := internalError.MapDBError(err)

		switch {
		case errors.Is(dbErr, internalError.ErrDBErrNotFound):
			return dto.AuthorResponse{}, internalError.AuthorNotFound

		default:
			s.logger.Error().Err(err).Msg("Unknown error")
			return dto.AuthorResponse{}, internalError.InternalError
		}
	}

	return dto.ToAuthorResponse(author), nil
}
