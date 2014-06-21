package nap

import (
	"net/http"
	"testing"
)

func TestCode(t *testing.T) {
	for _, test := range []struct {
		status Status
		code   int
	}{
		{status: OK{}, code: http.StatusOK},
		{status: Created{}, code: http.StatusCreated},
		{status: NotFound{}, code: http.StatusNotFound},
		{status: BadRequest{}, code: http.StatusBadRequest},
		{status: MethodNotAllowed{}, code: http.StatusMethodNotAllowed},
		{status: InternalError{}, code: http.StatusInternalServerError},
	} {
		if test.status.Code() != test.code {
			t.Fatalf("bad status (%q): got %d, want %d", test.status, test.status.Code(), test.code)
		}
	}
}

func TestMessage(t *testing.T) {
	for _, msg := range []string{
		"test message",
		"",
	} {
		for _, test := range []Status{
			OK{msg},
			Created{msg},
			NotFound{msg},
			BadRequest{msg},
			MethodNotAllowed{msg},
			InternalError{msg},
		} {
			if test.Message() != msg {
				t.Fatalf("bad status (%q): got %q, want %q", test, test.Message(), msg)
			}
		}
	}
}
