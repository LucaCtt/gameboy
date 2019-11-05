package assert

import (
	"reflect"
	"testing"
)

// Equal verifies that two values are deeply equal,
// as defined by the "reflect" package.
func Equal(t *testing.T, got, want interface{}) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

// Err verifies that the given error is (or is not) nil.
func Err(t *testing.T, got error, want bool) {
	t.Helper()
	if (got != nil) != want {
		t.Errorf("got %q, want %t", got, want)
	}
}
