package errors

import (
	"fmt"
	"testing"

	"github.com/lucactt/gameboy/util"
)

func TestError_Error(t *testing.T) {
	e := &Error{
		Message:   "test",
		Err:       fmt.Errorf("error"),
		Code:      CodeUnexpected,
		Component: Memory,
	}
	want := fmt.Sprintf("%s: %v", e.Message, e.Err)

	util.AssertEqual(t, e.Error(), want)
}

func TestError_Unwrap(t *testing.T) {
	e := &Error{
		Message:   "test",
		Err:       fmt.Errorf("error"),
		Code:      CodeUnexpected,
		Component: Memory,
	}
	want := e.Err

	util.AssertEqual(t, e.Unwrap(), want)
}
