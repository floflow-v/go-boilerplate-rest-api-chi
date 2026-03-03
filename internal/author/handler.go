package author

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/rs/zerolog"

	"go-boilerplate-rest-api-chi/internal/author/dto"
	internalError "go-boilerplate-rest-api-chi/internal/error"
	"go-boilerplate-rest-api-chi/internal/response"
	internalValidator "go-boilerplate-rest-api-chi/internal/validator"
)

type AuthorHandler struct {
	service   AuthorService
	validator *internalValidator.Validator
	logger    zerolog.Logger
}

func NewAuthorHandler(service AuthorService, validator *internalValidator.Validator, logger zerolog.Logger) *AuthorHandler {
	return &AuthorHandler{
		service:   service,
		validator: validator,
		logger:    logger,
	}
}

func (h *AuthorHandler) Routes() http.Handler {
	r := chi.NewRouter()

	// routes
	r.Post("/", h.CreateAuthor)
	r.Get("/{author_id}", h.GetAuthorByID)

	return r
}

// CreateAuthor godoc
//
//	@Summary		Create a new author
//	@Description	Create a new author with the provided data
//	@Tags			authors
//	@Accept			json
//	@Produce		json
//	@Param			author	body		dto.CreateAuthorRequest	true	"Author data"
//	@Success		201		{object}	dto.AuthorResponse
//	@Failure		400		{object}	internalError.Response
//	@Failure		409		{object}	internalError.Response
//	@Failure		500		{object}	internalError.Response
//	@Router			/authors [post]
func (h *AuthorHandler) CreateAuthor(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateAuthorRequest

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

	author, err := h.service.CreateAuthor(r.Context(), &req)
	if err != nil {
		internalError.Handle(w, err)
		return
	}

	authorResponse := dto.ToAuthorResponse(author)

	response.JSON(w, http.StatusCreated, authorResponse)
}

// GetAuthorByID godoc
//
//	@Summary		Get author by id
//	@Description	Get a single author by its ID
//	@Tags			authors
//	@Produce		json
//	@Param			author_id	path		string	true	"Author ID"
//	@Success		200			{object}	dto.AuthorResponse
//	@Failure		400			{object}	internalError.Response
//	@Failure		404			{object}	internalError.Response
//	@Failure		500			{object}	internalError.Response
//	@Router			/authors/{author_id} [get]
func (h *AuthorHandler) GetAuthorByID(w http.ResponseWriter, r *http.Request) {
	authorID, err := uuid.Parse(chi.URLParam(r, "author_id"))
	if err != nil {
		internalError.Handle(w, internalError.InvalidUUID)
		return
	}

	author, err := h.service.GetAuthorByID(r.Context(), authorID)
	if err != nil {
		internalError.Handle(w, err)
		return
	}

	authorResponse := dto.ToAuthorResponse(author)

	response.JSON(w, http.StatusOK, authorResponse)
}
