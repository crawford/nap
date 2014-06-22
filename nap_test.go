package nap

import (
	"net/http"
	"net/http/httptest"
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
		code    int
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
			code: http.StatusOK,
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
			code: http.StatusOK,
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
			code: http.StatusOK,
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
			code: http.StatusNotFound,
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
			code: http.StatusOK,
		},
	} {
		result, code := DefaultWrapper{}.Wrap(test.payload, test.status)
		if code != test.code {
			t.Fatalf("bad code (%q, %q): got %d, want %d", test.payload, test.status, code, test.code)
		}
		if !reflect.DeepEqual(result, test.result) {
			t.Fatalf("bad result (%q, %q): got %q, want %q", test.payload, test.status, result, test.result)
		}
	}
}

func TestHandlerFuncPanic(t *testing.T) {
	writer := httptest.NewRecorder()

	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("unexpected panic: %q", r)
		}
		if writer.Code != http.StatusInternalServerError {
			t.Fatalf("bad code: got %d, want %d", writer.Code, http.StatusInternalServerError)
		}
	}()

	HandlerFunc(func(*http.Request) (interface{}, Status) {
		panic("")
	}).ServeHTTP(writer, nil)
}

type PanicWrapper struct{}

func (w PanicWrapper) Wrap(payload interface{}, status Status) (interface{}, int) {
	panic("")
}

func TestHandlerFuncWrapperPanic(t *testing.T) {
	PayloadWrapper = PanicWrapper{}

	writer := httptest.NewRecorder()

	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("unexpected panic: %q", r)
		}
		if writer.Code != http.StatusInternalServerError {
			t.Fatalf("bad code: got %d, want %d", writer.Code, http.StatusInternalServerError)
		}
	}()

	HandlerFunc(func(*http.Request) (interface{}, Status) {
		return nil, nil
	}).ServeHTTP(writer, nil)
}
