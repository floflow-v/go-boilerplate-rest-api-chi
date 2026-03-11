package error

import (
	"errors"
	"net/http"

	"go-rest-api-chi-example/internal/response"
)

type Response struct {
	Error   string `json:"error" example:"error_code"`
	Message string `json:"message" example:"Human readable error message"`
	Details any    `json:"details,omitempty"`
}

func Handle(w http.ResponseWriter, err error) {
	var appErr *Error

	if errors.As(err, &appErr) {
		response.JSON(w, appErr.StatusCode, Response{
			Error:   appErr.Code,
			Message: appErr.Message,
			Details: appErr.Details,
		})
		return
	}

	response.JSON(w, http.StatusInternalServerError, Response{
		Error:   InternalError.Code,
		Message: InternalError.Message,
	})
}
