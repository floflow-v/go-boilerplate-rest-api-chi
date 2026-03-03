package error

type Error struct {
	Code       string
	Message    string
	StatusCode int
	Details    any
}

func (e *Error) Error() string {
	return e.Message
}

func New(code, message string, status int) *Error {
	return &Error{
		Code:       code,
		Message:    message,
		StatusCode: status,
	}
}

func (e *Error) WithDetails(details any) *Error {
	clone := *e
	clone.Details = details
	return &clone
}
