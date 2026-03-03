package book

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/rs/zerolog"

	"go-rest-api-chi-example/internal/book/dto"
	internalError "go-rest-api-chi-example/internal/error"
	"go-rest-api-chi-example/internal/response"
	internalValidator "go-rest-api-chi-example/internal/validator"
)

type BookHandler struct {
	service   BookService
	validator *internalValidator.Validator
	logger    zerolog.Logger
}

func NewBookHandler(service BookService, validator *internalValidator.Validator, logger zerolog.Logger) *BookHandler {
	return &BookHandler{
		service:   service,
		validator: validator,
		logger:    logger,
	}
}

func (h *BookHandler) Routes() http.Handler {
	r := chi.NewRouter()

	// routes
	r.Post("/", h.CreateBook)
	r.Get("/", h.GetAllBooks)
	r.Get("/{book_id}", h.GetBookByID)
	r.Patch("/{book_id}", h.UpdateBook)
	r.Delete("/{book_id}", h.DeleteBook)
	r.Get("/secure", h.AuthTestRoute)

	return r
}

// CreateBook godoc
//
//	@Summary		Create a new book
//	@Description	Create a new book with the provided data
//	@Tags			books
//	@Accept			json
//	@Produce		json
//	@Param			book	body		dto.CreateBookRequest	true	"Book data"
//	@Success		201		{object}	dto.BookResponse
//	@Failure		400		{object}	internalError.Response
//	@Failure		409		{object}	internalError.Response
//	@Failure		500		{object}	internalError.Response
//	@Router			/books [post]
func (h *BookHandler) CreateBook(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateBookRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		internalError.Handle(w, internalError.InvalidRequestBody)
		return
	}

	if err := h.validator.Struct(&req); err != nil {
		details := h.validator.FormatErrors(err)

		internalError.Handle(w,
			internalError.ValidationError.WithDetails(details),
		)
		return
	}

	book, err := h.service.CreateBook(r.Context(), &req)
	if err != nil {
		internalError.Handle(w, err)
		return
	}

	response.JSON(w, http.StatusCreated, book)
}

// GetAllBooks godoc
//
//	@Summary		Get all books
//	@Description	Get a list of all books
//	@Tags			books
//	@Produce		json
//	@Success		200	{object}	[]dto.BookResponse
//	@Failure		404	{object}	internalError.Response
//	@Failure		500	{object}	internalError.Response
//	@Router			/books [get]
func (h *BookHandler) GetAllBooks(w http.ResponseWriter, r *http.Request) {
	books, err := h.service.GetAllBooks(r.Context())
	if err != nil {
		internalError.Handle(w, err)
		return
	}

	response.JSON(w, http.StatusOK, books)
}

// GetBookByID godoc
//
//	@Summary		Get book by id
//	@Description	Get a single book by its ID
//	@Tags			books
//	@Produce		json
//	@Param			book_id	path		string	true	"Book ID"
//	@Success		200		{object}	dto.BookResponse
//	@Failure		400		{object}	internalError.Response
//	@Failure		404		{object}	internalError.Response
//	@Failure		500		{object}	internalError.Response
//	@Router			/books/{book_id} [get]
func (h *BookHandler) GetBookByID(w http.ResponseWriter, r *http.Request) {
	bookID, err := uuid.Parse(chi.URLParam(r, "book_id"))
	if err != nil {
		internalError.Handle(w, internalError.InvalidUUID)
		return
	}

	book, err := h.service.GetBookByID(r.Context(), bookID)
	if err != nil {
		internalError.Handle(w, err)
		return
	}

	response.JSON(w, http.StatusOK, book)
}

// UpdateBook godoc
//
//	@Summary		Update a book
//	@Description	Update a book with the provided data
//	@Tags			books
//	@Accept			json
//	@Produce		json
//	@Param			book_id	path		string					true	"Book ID"
//	@Param			book	body		dto.UpdateBookRequest	true	"Book data"
//	@Success		204		{object}	nil
//	@Failure		400		{object}	internalError.Response
//	@Failure		404		{object}	internalError.Response
//	@Failure		500		{object}	internalError.Response
//	@Router			/books/{book_id} [patch]
func (h *BookHandler) UpdateBook(w http.ResponseWriter, r *http.Request) {
	bookID, err := uuid.Parse(chi.URLParam(r, "book_id"))
	if err != nil {
		internalError.Handle(w, internalError.InvalidUUID)
		return
	}

	var req dto.UpdateBookRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		internalError.Handle(w, internalError.InvalidRequestBody)
		return
	}

	if err := h.validator.Struct(&req); err != nil {
		details := h.validator.FormatErrors(err)
		internalError.Handle(w,
			internalError.ValidationError.WithDetails(details),
		)
		return
	}

	err = h.service.UpdateBook(r.Context(), &req, bookID)
	if err != nil {
		internalError.Handle(w, err)
		return
	}

	response.NoContent(w)
}

// DeleteBook godoc
//
//	@Summary		Delete a book
//	@Description	Delete a book by its ID
//	@Tags			books
//	@Produce		json
//	@Param			book_id	path		string	true	"Book ID"
//	@Success		204		{object}	nil
//	@Failure		404		{object}	internalError.Response
//	@Failure		500		{object}	internalError.Response
//	@Router			/books/{book_id} [delete]
func (h *BookHandler) DeleteBook(w http.ResponseWriter, r *http.Request) {
	bookID, err := uuid.Parse(chi.URLParam(r, "book_id"))
	if err != nil {
		internalError.Handle(w, internalError.InvalidUUID)
		return
	}

	err = h.service.DeleteBook(r.Context(), bookID)
	if err != nil {
		internalError.Handle(w, err)
		return
	}

	response.NoContent(w)
}

// AuthTestRoute godoc
//
//	@Summary		Authenticated test route
//	@Description	Authenticated test route
//	@Tags			books
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Success		200	{object}	response.SuccessResponse
//	@Router			/books/secure [get]
func (h *BookHandler) AuthTestRoute(w http.ResponseWriter, _ *http.Request) {
	response.JSON(w, http.StatusOK, response.SuccessResponse{
		Status:  "success",
		Message: "ok",
	})
}
