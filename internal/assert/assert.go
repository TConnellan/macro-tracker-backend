package assert

import (
	"errors"
	"strings"
	"testing"
)

func Equal[T comparable](t *testing.T, actual, expected T) {
	t.Helper()

	if actual != expected {
		t.Errorf("got: %v; want: %v", actual, expected)
	}
}

func NotEqual[T comparable](t *testing.T, actual, expected T) {
	t.Helper()

	if actual == expected {
		t.Errorf("got: %v; want: %v", actual, expected)
	}
}

func StringContains(t *testing.T, actual, expectedSubstring string) {
	t.Helper()

	if !strings.Contains(actual, expectedSubstring) {
		t.Errorf("got: %q; expected to contain: %q", actual, expectedSubstring)
	}
}

func NilError(t *testing.T, actual error) {
	t.Helper()
	if actual != nil {
		t.Errorf("got: %v; expected: nil", actual)
	}
}

func ExpectError(t *testing.T, err error, expected error) {
	t.Helper()
	if !errors.Is(err, expected) {
		t.Errorf("got: %v; expected: %v", err, expected)
	}
}

func FailWithError(t *testing.T, err error) {
	t.Helper()
	t.Errorf("encountered error: %v", err)
}
