package nap

import (
	"encoding/json"
	"net/http"
)

var (
	MethodNotAllowedHandler HandlerFunc  = HandlerFunc(defaultMethodNotAllowed)
	NotFoundHandler         HandlerFunc  = HandlerFunc(defaultNotFound)
	PayloadWrapper          Wrapper      = DefaultWrapper{}
	PanicHandler            ErrorHandler = nil
	ResponseHeaders         []Header     = defaultHeaders
)

type Wrapper interface {
	Wrap(payload interface{}, status Status) (interface{}, int)
}

type ErrorHandler interface {
	Handle(e interface{})
}

type DefaultWrapper struct{}

func (w DefaultWrapper) Wrap(payload interface{}, status Status) (interface{}, int) {
	if status == nil {
		status = OK{}
	}
	return payload, status.Code()
}

type Header struct {
	Name  string
	Value []string
}

type HandlerFunc func(req *http.Request) (interface{}, Status)

func (f HandlerFunc) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	var payload interface{}
	var status Status
	defer func() {
		var result interface{}
		var code int
		defer func() {
			for _, header := range ResponseHeaders {
				writer.Header()[header.Name] = header.Value
			}

			if r := recover(); r != nil {
				handlePanic(r)
			} else {
				if res, err := json.Marshal(result); err == nil {
					writer.WriteHeader(code)
					writer.Write(res)
					return
				}
			}
			writer.WriteHeader(http.StatusInternalServerError)
		}()
		if r := recover(); r != nil {
			handlePanic(r)
			status = InternalError{}
		}
		result, code = PayloadWrapper.Wrap(payload, status)
	}()
	payload, status = f(request)
}

func handlePanic(r interface{}) {
	if PanicHandler != nil {
		func() {
			defer func() {
				recover()
			}()
			PanicHandler.Handle(r)
		}()
	}
}

var (
	defaultHeaders = []Header{
		{"Content-Type", []string{"application/json"}},
		{"Cache-Control", []string{"no-cache,must-revalidate"}},
	}
)

func defaultMethodNotAllowed(req *http.Request) (interface{}, Status) {
	return nil, MethodNotAllowed{"method not allowed on resource"}
}

func defaultNotFound(req *http.Request) (interface{}, Status) {
	return nil, NotFound{"resource not found"}
}
