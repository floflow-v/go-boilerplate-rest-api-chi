package author

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"gorm.io/gorm"

	internalError "go-boilerplate-rest-api-chi/internal/error"
	"go-boilerplate-rest-api-chi/internal/model"
)

//go:generate mockgen -destination=../mocks/mock_author_repository.go -package=mocks go-boilerplate-rest-api-chi/internal/author AuthorRepository
type AuthorRepository interface {
	Create(ctx context.Context, newAuthor *model.Author) (*model.Author, error)
	GetByID(ctx context.Context, authorID uuid.UUID) (*model.Author, error)
	Exists(ctx context.Context, authorID uuid.UUID) (bool, error)
}

type authorRepository struct {
	db     *gorm.DB
	logger zerolog.Logger
}

func NewAuthorRepository(db *gorm.DB, logger zerolog.Logger) AuthorRepository {
	return &authorRepository{
		db:     db,
		logger: logger,
	}
}

func (r *authorRepository) Create(ctx context.Context, newAuthor *model.Author) (*model.Author, error) {
	if err := r.db.WithContext(ctx).Create(newAuthor).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, internalError.AuthorDuplicate
		}

		r.logger.Error().Err(err).Msg("database error")
		return nil, err
	}

	return newAuthor, nil
}

func (r *authorRepository) GetByID(ctx context.Context, authorID uuid.UUID) (*model.Author, error) {
	var author *model.Author

	if err := r.db.WithContext(ctx).First(&author, "id = ?", authorID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, internalError.AuthorNotFound
		}

		r.logger.Error().Err(err).Msg("database error")
		return nil, err
	}

	return author, nil
}

func (r *authorRepository) Exists(ctx context.Context, authorID uuid.UUID) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.Author{}).Where("id = ?", authorID).Count(&count).Error
	return count > 0, err
}
