package author

import (
	"context"

	"github.com/google/uuid"
	"github.com/rs/zerolog"

	"go-boilerplate-rest-api-chi/internal/author/dto"
	"go-boilerplate-rest-api-chi/internal/model"
)

//go:generate mockgen -destination=../mocks/mock_author_service.go -package=mocks go-boilerplate-rest-api-chi/internal/author AuthorService
type AuthorService interface {
	CreateAuthor(ctx context.Context, req *dto.CreateAuthorRequest) (*model.Author, error)
	GetAuthorByID(ctx context.Context, authorID uuid.UUID) (*model.Author, error)
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

func (s *authorService) CreateAuthor(ctx context.Context, req *dto.CreateAuthorRequest) (*model.Author, error) {
	author := &model.Author{
		Name: req.Name,
	}

	return s.repository.Create(ctx, author)
}

func (s *authorService) GetAuthorByID(ctx context.Context, authorID uuid.UUID) (*model.Author, error) {
	author, err := s.repository.GetByID(ctx, authorID)
	if err != nil {
		return nil, err
	}

	return author, nil
}
