package errors

import (
	"errors"
	"net/http"
	"testing"
)

func TestStatusCodeReturnsAppErrorStatus(t *testing.T) {
	err := InvalidInput("bad input", nil)
	if got := StatusCode(err); got != http.StatusBadRequest {
		t.Fatalf("expected %d, got %d", http.StatusBadRequest, got)
	}
}

func TestStatusCodeWrapsUnknownErrors(t *testing.T) {
	err := errors.New("boom")
	if got := StatusCode(err); got != http.StatusInternalServerError {
		t.Fatalf("expected %d, got %d", http.StatusInternalServerError, got)
	}
}
