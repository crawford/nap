package nap

import (
	"net/http"
	"testing"
)

func TestDefaults(t *testing.T) {
	for _, test := range []struct {
		fn      func(req *http.Request) (interface{}, Status)
		payload interface{}
		status  Status
	}{
		{
			fn:      defaultMethodNotAllowed,
			payload: nil,
			status:  MethodNotAllowed{"method not allowed on resource"},
		},
		{
			fn:      defaultNotFound,
			payload: nil,
			status:  NotFound{"resource not found"},
		},
	} {
		payload, status := test.fn(nil)
		if payload != test.payload {
			t.Fatalf("bad payload: got %q, want %q", payload, test.payload)
		}
		if status != test.status {
			t.Fatalf("bad status: got %q, want %q", status, test.status)
		}
	}
}
