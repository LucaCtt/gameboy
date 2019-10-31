package errors

import (
	"fmt"
	"testing"

	"github.com/lucactt/gameboy/util/assert"
)

func Test_E(t *testing.T) {
	t.Run("valid args", func(t *testing.T) {
		got := E("test", fmt.Errorf("test"), CodeUnexpected, Memory)
		want := &Error{
			Message:   "test",
			Err:       fmt.Errorf("test"),
			Code:      CodeUnexpected,
			Component: Memory,
		}

		assert.Equal(t, got, want)
	})

	t.Run("invalid arg", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("did not panic")
			}
		}()
		E(69.420)
	})
}

func Test_Code(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want ErrCode
	}{
		{"err type is not *Error", fmt.Errorf("test"), CodeUnexpected},
		{"err code is not 0", &Error{Code: CodeOutOfRange, Err: fmt.Errorf("test")}, CodeOutOfRange},
		{"wrapped err is *Error", &Error{Err: &Error{Code: CodeOutOfRange}}, CodeOutOfRange},
		{"wrapped err is not *Error", &Error{Err: fmt.Errorf("test")}, CodeUnexpected},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Code(tt.err)

			if got != tt.want {
				t.Errorf("got %d, want %d", got, tt.want)
			}
		})
	}
}
