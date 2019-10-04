package util

import (
	"reflect"
	"testing"
)

// AssertEqual verifies that two values are deeply equal,
// as defined by the "reflect" package.
func AssertEqual(t *testing.T, got, want interface{}) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

// AssertErr verifies that the given error is (or is not) nil.
func AssertErr(t *testing.T, got error, want bool) {
	t.Helper()
	AssertEqual(t, got != nil, want)
}
