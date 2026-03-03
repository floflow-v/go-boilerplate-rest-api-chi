package book

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"gorm.io/gorm"

	internalError "go-boilerplate-rest-api-chi/internal/error"
	"go-boilerplate-rest-api-chi/internal/model"
)

//go:generate mockgen -destination=../mocks/mock_book_repository.go -package=mocks go-boilerplate-rest-api-chi/internal/book BookRepository
type BookRepository interface {
	Create(ctx context.Context, book *model.Book) (*model.Book, error)
	GetAll(ctx context.Context) ([]*model.Book, error)
	GetByID(ctx context.Context, bookID uuid.UUID) (*model.Book, error)
	Update(ctx context.Context, bookID uuid.UUID, updates map[string]interface{}) error
	Delete(ctx context.Context, bookID uuid.UUID) error
}

type bookRepository struct {
	db     *gorm.DB
	logger zerolog.Logger
}

func NewBookRepository(db *gorm.DB, logger zerolog.Logger) BookRepository {
	return &bookRepository{
		db:     db,
		logger: logger,
	}
}

func (r *bookRepository) Create(ctx context.Context, newBook *model.Book) (*model.Book, error) {
	if err := r.db.WithContext(ctx).Create(newBook).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			r.logger.Error().Err(err).Msg("record already exist in database")
			return nil, internalError.BookDuplicate
		}

		r.logger.Error().Err(err).Msg("database error")
		return nil, err
	}

	return newBook, nil
}

func (r *bookRepository) GetAll(ctx context.Context) ([]*model.Book, error) {
	var books []*model.Book

	if err := r.db.WithContext(ctx).Preload("Author").Find(&books).Error; err != nil {
		r.logger.Error().Err(err).Msg("error when retrieve books on database ")
		return nil, err
	}

	return books, nil
}

func (r *bookRepository) GetByID(ctx context.Context, bookID uuid.UUID) (*model.Book, error) {
	book := &model.Book{}

	if err := r.db.WithContext(ctx).Preload("Author").First(&book, "id = ?", bookID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, internalError.BookNotFound
		}
		return nil, err
	}

	return book, nil
}

func (r *bookRepository) Update(ctx context.Context, bookID uuid.UUID, updates map[string]interface{}) error {
	result := r.db.WithContext(ctx).Model(&model.Book{ID: bookID}).Updates(updates)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return internalError.BookNotFound
	}

	return nil
}

func (r *bookRepository) Delete(ctx context.Context, bookID uuid.UUID) error {
	result := r.db.WithContext(ctx).Delete(&model.Book{ID: bookID})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return internalError.BookNotFound
	}

	return nil
}
