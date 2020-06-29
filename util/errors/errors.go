package errors

import (
	"fmt"
)

// ErrCode identifies the error with a code.
type ErrCode int

// ErrComponent identifies the component where the error happened.
type ErrComponent string

// Components where errors can be originated from.
const (
	Mem  ErrComponent = "memory"
	Cart ErrComponent = "cartridge"
	CPU  ErrComponent = "CPU"
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
