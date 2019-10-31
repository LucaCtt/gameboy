package errors

import (
	"fmt"
	"testing"

	"github.com/lucactt/gameboy/util/assert"
)

func TestError_Error(t *testing.T) {
	e := &Error{
		Message:   "test",
		Err:       fmt.Errorf("error"),
		Code:      CodeUnexpected,
		Component: Memory,
	}
	want := fmt.Sprintf("%s: %v", e.Message, e.Err)

	assert.Equal(t, e.Error(), want)
}

func TestError_Unwrap(t *testing.T) {
	e := &Error{
		Message:   "test",
		Err:       fmt.Errorf("error"),
		Code:      CodeUnexpected,
		Component: Memory,
	}
	want := e.Err

	assert.Equal(t, e.Unwrap(), want)
}
