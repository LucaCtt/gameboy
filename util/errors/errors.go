package errors

import (
	"fmt"
	"net/http"
)

type ErrCode int
type ErrComponent string

// HTTP status codes used to identify errors.
const (
	CodeUnexpected ErrCode = http.StatusInternalServerError
	CodeOutOfRange ErrCode = http.StatusInsufficientStorage
)

// Components where errors can be originated from.
const (
	Memory    ErrComponent = "memory"
	Cartridge ErrComponent = "cartridge"
)

// Error is a wrapper for an error value with added context.
type Error struct {
	Message   string
	Err       error
	Code      ErrCode
	Component ErrComponent
}

func (e *Error) Error() string {
	return fmt.Sprintf("%s: %v", e.Message, e.Err)
}

// Unwrap returns the wrapped error.
func (e *Error) Unwrap() error {
	return e.Err
}
