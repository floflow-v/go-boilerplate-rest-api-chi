package author

import (
	"context"

	"github.com/google/uuid"
	"github.com/rs/zerolog"

	"go-rest-api-chi-example/internal/author/dto"
	"go-rest-api-chi-example/internal/model"
)

//go:generate mockgen -destination=../mocks/mock_author_service.go -package=mocks go-rest-api-chi-example/internal/author AuthorService
type AuthorService interface {
	CreateAuthor(ctx context.Context, req *dto.CreateAuthorRequest) (*dto.AuthorResponse, error)
	GetAuthorByID(ctx context.Context, authorID uuid.UUID) (*dto.AuthorResponse, error)
}

type authorService struct {
	repository AuthorRepository
	logger     zerolog.Logger
}

func NewAuthorService(repository AuthorRepository, logger zerolog.Logger) AuthorService {
	return &authorService{
		repository: repository,
		logger:     logger,
	}
}

func (s *authorService) CreateAuthor(ctx context.Context, req *dto.CreateAuthorRequest) (*dto.AuthorResponse, error) {
	author := &model.Author{
		Name: req.Name,
	}

	newAuthor, err := s.repository.Create(ctx, author)
	if err != nil {
		return nil, err
	}

	return dto.ToAuthorResponse(newAuthor), nil
}

func (s *authorService) GetAuthorByID(ctx context.Context, authorID uuid.UUID) (*dto.AuthorResponse, error) {
	author, err := s.repository.GetByID(ctx, authorID)
	if err != nil {
		return nil, err
	}

	return dto.ToAuthorResponse(author), nil
}
