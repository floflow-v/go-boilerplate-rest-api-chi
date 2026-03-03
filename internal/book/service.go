package book

import (
	"context"

	"github.com/google/uuid"
	"github.com/rs/zerolog"

	"go-rest-api-chi-example/internal/author"
	"go-rest-api-chi-example/internal/book/dto"
	internalError "go-rest-api-chi-example/internal/error"
	"go-rest-api-chi-example/internal/model"
)

//go:generate mockgen -destination=../mocks/mock_book_service.go -package=mocks go-rest-api-chi-example/internal/book BookService
type BookService interface {
	CreateBook(ctx context.Context, req *dto.CreateBookRequest) (*dto.BookResponse, error)
	GetAllBooks(ctx context.Context) ([]dto.BookResponse, error)
	GetBookByID(ctx context.Context, bookID uuid.UUID) (*dto.BookResponse, error)
	UpdateBook(ctx context.Context, req *dto.UpdateBookRequest, bookID uuid.UUID) error
	DeleteBook(ctx context.Context, bookID uuid.UUID) error
}

type bookService struct {
	repository       BookRepository
	authorRepository author.AuthorRepository
	logger           zerolog.Logger
}

func NewBookService(repository BookRepository, authorRepository author.AuthorRepository, logger zerolog.Logger) BookService {
	return &bookService{
		repository:       repository,
		authorRepository: authorRepository,
		logger:           logger,
	}
}

func (s *bookService) CreateBook(ctx context.Context, req *dto.CreateBookRequest) (*dto.BookResponse, error) {
	authorID, err := uuid.Parse(req.AuthorID)
	if err != nil {
		return nil, internalError.InvalidAuthorID
	}

	exists, err := s.authorRepository.Exists(ctx, authorID)
	if err != nil {
		return nil, err
	}

	if !exists {
		return nil, internalError.AuthorNotFound
	}

	book := &model.Book{
		Title:       req.Title,
		Description: req.Description,
		AuthorID:    authorID,
	}

	newBook, err := s.repository.Create(ctx, book)
	if err != nil {
		return nil, err
	}

	fetchedBook, err := s.repository.GetByID(ctx, newBook.ID)
	if err != nil {
		return nil, err
	}

	return dto.ToBookResponse(fetchedBook), nil
}

func (s *bookService) GetAllBooks(ctx context.Context) ([]dto.BookResponse, error) {
	books, err := s.repository.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	if len(books) == 0 {
		return nil, internalError.BookNotFound
	}

	return dto.ToBooksResponse(books), nil
}

func (s *bookService) GetBookByID(ctx context.Context, bookID uuid.UUID) (*dto.BookResponse, error) {
	book, err := s.repository.GetByID(ctx, bookID)
	if err != nil {
		return nil, err
	}

	return dto.ToBookResponse(book), nil
}

func (s *bookService) UpdateBook(ctx context.Context, req *dto.UpdateBookRequest, bookID uuid.UUID) error {
	updates := map[string]interface{}{
		"description": req.Description,
	}

	return s.repository.Update(ctx, bookID, updates)
}

func (s *bookService) DeleteBook(ctx context.Context, bookID uuid.UUID) error {
	return s.repository.Delete(ctx, bookID)
}
