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
			result:  nil,
			code:    http.StatusOK,
		},
		{
			payload: nil,
			status:  OK{},
			result:  nil,
			code:    http.StatusOK,
		},
		{
			payload: nil,
			status:  OK{"test"},
			result:  nil,
			code:    http.StatusOK,
		},
		{
			payload: nil,
			status:  NotFound{},
			result:  nil,
			code:    http.StatusNotFound,
		},
		{
			payload: "test",
			status:  nil,
			result:  "test",
			code:    http.StatusOK,
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

func TestHeaders(t *testing.T) {
	for _, test := range []struct {
		headers []Header
	}{
		{
			headers: nil,
		},
		{
			headers: []Header{},
		},
		{
			headers: defaultHeaders,
		},
	} {
		handler := HandlerFunc(func(*http.Request) (interface{}, Status) {
			return nil, nil
		})
		func() {
			writer := httptest.NewRecorder()
			defer func() {
				if r := recover(); r != nil {
					t.Fatalf("unexpected panic: %q", r)
				}
				for _, header := range test.headers {
					if !reflect.DeepEqual(writer.Header()[header.Name], header.Value) {
						t.Fatalf("bad header: got %q, want %q", writer.Header()[header.Name], header.Value)
					}
					delete(writer.Header(), header.Name)
				}
				if len(writer.Header()) > 0 {
					t.Fatalf("extra headers: got %d, want %d", len(writer.Header()), 0)
				}
			}()
			ResponseHeaders = test.headers
			handler.ServeHTTP(writer, nil)
		}()
	}
}

type PanicWrapper struct{}

func (w PanicWrapper) Wrap(payload interface{}, status Status) (interface{}, int) {
	panic("")
}

func TestHandlerFunc(t *testing.T) {
	for _, test := range []struct {
		wrapper Wrapper
		handler HandlerFunc
		code    int
		body    string
	}{
		{
			wrapper: DefaultWrapper{},
			handler: HandlerFunc(func(*http.Request) (interface{}, Status) {
				panic("")
			}),
			code: http.StatusInternalServerError,
			body: `null`,
		},
		{
			wrapper: PanicWrapper{},
			handler: HandlerFunc(func(*http.Request) (interface{}, Status) {
				return nil, nil
			}),
			code: http.StatusInternalServerError,
			body: ``,
		},
		{
			wrapper: DefaultWrapper{},
			handler: HandlerFunc(func(*http.Request) (interface{}, Status) {
				return nil, NotFound{}
			}),
			code: http.StatusNotFound,
			body: `null`,
		},
		{
			wrapper: DefaultWrapper{},
			handler: HandlerFunc(func(*http.Request) (interface{}, Status) {
				return "testing", OK{"test"}
			}),
			code: http.StatusOK,
			body: `"testing"`,
		},
	} {
		func() {
			writer := httptest.NewRecorder()
			defer func() {
				if r := recover(); r != nil {
					t.Fatalf("unexpected panic: %q", r)
				}
				if writer.Code != test.code {
					t.Fatalf("bad code: got %d, want %d", writer.Code, test.code)
				}
				if writer.Body.String() != test.body {
					t.Fatalf("bad body: got %q, want %q", writer.Body.String(), test.body)
				}
			}()
			PayloadWrapper = test.wrapper
			test.handler.ServeHTTP(writer, nil)
		}()
	}
}
