package errors

import (
	"errors"
	"testing"
)

func TestNewf(t *testing.T) {
	err := Newf("error: %s %d", "test", 42)
	if err == nil {
		t.Fatal("expected non-nil error")
	}
	expected := "error: test 42"
	if !errors.Is(err, err) || err.Error() != expected {
		t.Errorf("expected error message %q, got %q", expected, err.Error())
	}
}
