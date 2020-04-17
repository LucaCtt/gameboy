package errors

import (
	"fmt"
	"net/http"
)

// ErrCode identifies the error with a code.
type ErrCode int

// ErrComponent identifies the component where the error happened.
type ErrComponent string

// HTTP status codes used to identify errors.
const (
	CodeUnexpected ErrCode = http.StatusInternalServerError
	CodeOutOfRange ErrCode = http.StatusInsufficientStorage
)

// Components where errors can be originated from.
const (
	Mem  ErrComponent = "memory"
	Cart ErrComponent = "cartridge"
)

// Error is a wrapper for an error value with added context.
type Error struct {
	Message   string
	Err       error
	Code      ErrCode
	Component ErrComponent
}

func (e *Error) Error() string {
	if e.Err == nil {
		return fmt.Sprintf("%s", e.Message)
	}
	return fmt.Sprintf("%s: %v", e.Message, e.Err)
}

// Unwrap returns the wrapped error.
func (e *Error) Unwrap() error {
	return e.Err
}
