package nap

import (
	"encoding/json"
	"net/http"
)

var (
	MethodNotAllowedHandler HandlerFunc = HandlerFunc(defaultMethodNotAllowed)
	NotFoundHandler         HandlerFunc = HandlerFunc(defaultNotFound)
	PayloadWrapper          Wrapper     = DefaultWrapper{}
)

type Wrapper interface {
	Wrap(payload interface{}, status Status) (interface{}, int)
}

type DefaultWrapper struct{}

func (w DefaultWrapper) Wrap(payload interface{}, status Status) (interface{}, int) {
	if status == nil {
		status = OK{}
	}
	return map[string]interface{}{
		"status": map[string]interface{}{
			"code":    status.Code(),
			"message": status.Message(),
		},
		"result": payload,
	}, status.Code()

}

type HandlerFunc func(req *http.Request) (interface{}, Status)

func (f HandlerFunc) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	var payload interface{}
	var status Status
	defer func() {
		var result interface{}
		var code int
		defer func() {
			writer.Header().Add("Content-Type", "application/json")
			writer.Header().Add("Cache-Control", "no-cache,must-revalidate")

			if r := recover(); r == nil {
				if res, err := json.Marshal(result); err == nil {
					writer.WriteHeader(code)
					writer.Write(res)
					return
				}
			}
			writer.WriteHeader(http.StatusInternalServerError)
		}()
		if r := recover(); r != nil {
			status = InternalError{}
		}
		result, code = PayloadWrapper.Wrap(payload, status)
	}()
	payload, status = f(request)
}

func defaultMethodNotAllowed(req *http.Request) (interface{}, Status) {
	return nil, MethodNotAllowed{"method not allowed on resource"}
}

func defaultNotFound(req *http.Request) (interface{}, Status) {
	return nil, NotFound{"resource not found"}
}
