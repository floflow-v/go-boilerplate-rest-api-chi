package book

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/rs/zerolog"

	"go-rest-api-chi-example/internal/book/dto"
	"go-rest-api-chi-example/internal/database/sqlc"
	internalError "go-rest-api-chi-example/internal/error"
)

//go:generate mockgen -destination=../mocks/mock_book_service.go -package=mocks go-rest-api-chi-example/internal/book BookService
type BookService interface {
	CreateBook(ctx context.Context, req dto.CreateBookRequest) (dto.BookResponse, error)
	GetAllBooks(ctx context.Context) ([]dto.BookResponse, error)
	GetBookByID(ctx context.Context, bookID uuid.UUID) (dto.BookResponse, error)
	UpdateBook(ctx context.Context, req dto.UpdateBookRequest, bookID uuid.UUID) error
	DeleteBook(ctx context.Context, bookID uuid.UUID) error
}

type bookService struct {
	querier sqlc.Querier
	logger  zerolog.Logger
}

func NewBookService(querier sqlc.Querier, logger zerolog.Logger) BookService {
	return &bookService{
		querier: querier,
		logger:  logger,
	}
}

func (s *bookService) CreateBook(ctx context.Context, req dto.CreateBookRequest) (dto.BookResponse, error) {
	authorID, err := uuid.Parse(req.AuthorID)
	if err != nil {
		return dto.BookResponse{}, internalError.InvalidAuthorID
	}

	_, err = s.querier.GetAuthorByID(ctx, authorID.String())
	if err != nil {
		dbErr := internalError.MapDBError(err)

		switch {
		case errors.Is(dbErr, internalError.ErrDBErrNotFound):
			return dto.BookResponse{}, internalError.AuthorNotFound

		default:
			s.logger.Error().Err(err).Msg("Unknown error")
			return dto.BookResponse{}, internalError.InternalError
		}
	}

	bookID, _ := uuid.NewV7()

	bookRequest := sqlc.CreateBookParams{
		ID:          bookID.String(),
		Title:       req.Title,
		Description: req.Description,
		AuthorID:    authorID.String(),
	}

	err = s.querier.CreateBook(ctx, bookRequest)
	if err != nil {
		dbErr := internalError.MapDBError(err)

		switch {
		case errors.Is(dbErr, internalError.ErrDBErrDuplicate):
			return dto.BookResponse{}, internalError.BookDuplicate

		default:
			s.logger.Error().Err(err).Msg("Unknown error")
			return dto.BookResponse{}, internalError.InternalError
		}
	}

	book, err := s.querier.GetBookByID(ctx, bookID.String())
	if err != nil {
		s.logger.Error().Err(err).Msg("failed to get book after creation")
		return dto.BookResponse{}, internalError.InternalError
	}

	return dto.ToBookResponse(book), nil
}

func (s *bookService) GetAllBooks(ctx context.Context) ([]dto.BookResponse, error) {
	books, err := s.querier.GetAllBooks(ctx)
	if err != nil {
		return nil, internalError.InternalError
	}

	if len(books) == 0 {
		return nil, internalError.BookNotFound
	}

	return dto.ToBookResponseFromRows(books), nil
}

func (s *bookService) GetBookByID(ctx context.Context, bookID uuid.UUID) (dto.BookResponse, error) {
	book, err := s.querier.GetBookByID(ctx, bookID.String())
	if err != nil {
		dbErr := internalError.MapDBError(err)

		switch {
		case errors.Is(dbErr, internalError.ErrDBErrNotFound):
			return dto.BookResponse{}, internalError.BookNotFound
		default:
			s.logger.Error().Err(err).Msg("Unknown error")
			return dto.BookResponse{}, internalError.InternalError
		}
	}

	return dto.ToBookResponse(book), nil
}

func (s *bookService) UpdateBook(ctx context.Context, req dto.UpdateBookRequest, bookID uuid.UUID) error {
	_, err := s.querier.GetBookByID(ctx, bookID.String())
	if err != nil {
		dbErr := internalError.MapDBError(err)

		switch {
		case errors.Is(dbErr, internalError.ErrDBErrNotFound):
			return internalError.BookNotFound
		default:
			s.logger.Error().Err(err).Msg("Unknown error")
			return internalError.InternalError
		}
	}

	bookRequest := sqlc.UpdateBookParams{
		Title:       req.Title,
		Description: req.Description,
		ID:          bookID.String(),
	}

	err = s.querier.UpdateBook(ctx, bookRequest)
	if err != nil {
		s.logger.Error().Err(err).Msg("Unknown error")
		return internalError.InternalError
	}

	return nil
}

func (s *bookService) DeleteBook(ctx context.Context, bookID uuid.UUID) error {
	_, err := s.querier.GetBookByID(ctx, bookID.String())
	if err != nil {
		dbErr := internalError.MapDBError(err)

		switch {
		case errors.Is(dbErr, internalError.ErrDBErrNotFound):
			return internalError.BookNotFound
		default:
			s.logger.Error().Err(err).Msg("Unknown error")
			return internalError.InternalError
		}
	}

	err = s.querier.DeleteBook(ctx, bookID.String())
	if err != nil {
		s.logger.Error().Err(err).Msg("Unknown error")
		return internalError.InternalError
	}

	return nil
}
