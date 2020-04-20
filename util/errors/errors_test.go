package errors

import (
	"fmt"
	"testing"

	"github.com/lucactt/gameboy/util/assert"
)

func TestError_Error(t *testing.T) {
	t.Run("No wrapper err", func(t *testing.T) {
		e := &Error{
			Message:   "test",
			Err:       nil,
			Component: Mem,
		}

		assert.Equal(t, e.Error(), "test")
	})

	t.Run("Wrapped err", func(t *testing.T) {
		e := &Error{
			Message:   "test",
			Err:       fmt.Errorf("error"),
			Component: Mem,
		}
		want := fmt.Sprintf("%s: %v", e.Message, e.Err)

		assert.Equal(t, e.Error(), want)
	})
}

func TestError_Unwrap(t *testing.T) {
	e := &Error{
		Message:   "test",
		Err:       fmt.Errorf("error"),
		Component: Mem,
	}
	want := e.Err

	assert.Equal(t, e.Unwrap(), want)
}
