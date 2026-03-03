package error

import "net/http"

var (
	// authors
	AuthorNotFound  = New("author_not_found", "Author not found", http.StatusNotFound)
	AuthorDuplicate = New("author_duplicate", "Author already exists", http.StatusConflict)

	// books
	BookNotFound    = New("book_not_found", "Book not found", http.StatusNotFound)
	BookDuplicate   = New("book_duplicate", "Book already exists", http.StatusConflict)
	InvalidAuthorID = New("invalid_author_id", "Invalid author ID", http.StatusBadRequest)

	// general
	ValidationError    = New("validation_error", "Validation failed", http.StatusBadRequest)
	InvalidRequestBody = New("invalid_request_body", "Invalid request body", http.StatusBadRequest)
	InvalidUUID        = New("invalid_uuid", "Invalid uuid", http.StatusBadRequest)

	InternalError = New("internal_error", "Internal server error", http.StatusInternalServerError)
)
