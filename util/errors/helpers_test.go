package errors

import (
	"fmt"
	"testing"

	"github.com/lucactt/gameboy/util/assert"
)

func Test_E(t *testing.T) {
	t.Run("valid args", func(t *testing.T) {
		got := E("test", fmt.Errorf("test"), Mem)
		want := &Error{
			Message:   "test",
			Err:       fmt.Errorf("test"),
			Component: Mem,
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
