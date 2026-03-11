package error

import (
	"database/sql"
	"errors"

	"github.com/go-sql-driver/mysql"
)

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

func MapDBError(err error) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, sql.ErrNoRows) {
		return ErrDBErrNotFound
	}

	var mysqlErr *mysql.MySQLError

	if errors.As(err, &mysqlErr) {
		switch mysqlErr.Number {
		case 1062:
			return ErrDBErrDuplicate
		case 1451, 1452:
			return ErrDBErrForeignKeyViolation
		case 1048:
			return ErrDBErrNotNullViolation
		case 1406:
			return ErrDBErrDataTooLong
		}
	}
	return err
}
