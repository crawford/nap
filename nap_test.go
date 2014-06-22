package nap

import (
	"net/http"
	"reflect"
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
			t.Fatalf("bad payload (%q): got %q, want %q", test.fn, payload, test.payload)
		}
		if status != test.status {
			t.Fatalf("bad status (%q): got %q, want %q", test.fn, status, test.status)
		}
	}
}

func TestDefaultWrapper(t *testing.T) {
	for _, test := range []struct {
		payload interface{}
		status  Status
		result  interface{}
	}{
		{
			payload: nil,
			status:  nil,
			result: map[string]interface{}{
				"status": map[string]interface{}{
					"code":    http.StatusOK,
					"message": "",
				},
				"result": nil,
			},
		},
		{
			payload: nil,
			status:  OK{},
			result: map[string]interface{}{
				"status": map[string]interface{}{
					"code":    http.StatusOK,
					"message": "",
				},
				"result": nil,
			},
		},
		{
			payload: nil,
			status:  OK{"test"},
			result: map[string]interface{}{
				"status": map[string]interface{}{
					"code":    http.StatusOK,
					"message": "test",
				},
				"result": nil,
			},
		},
		{
			payload: nil,
			status:  NotFound{},
			result: map[string]interface{}{
				"status": map[string]interface{}{
					"code":    http.StatusNotFound,
					"message": "",
				},
				"result": nil,
			},
		},
		{
			payload: "test",
			status:  nil,
			result: map[string]interface{}{
				"status": map[string]interface{}{
					"code":    http.StatusOK,
					"message": "",
				},
				"result": "test",
			},
		},
	} {
		result := DefaultWrapper{}.Wrap(test.payload, test.status)
		if !reflect.DeepEqual(result, test.result) {
			t.Fatalf("bad result (%q, %q): got %q, want %q", test.payload, test.status, result, test.result)
		}
	}
}