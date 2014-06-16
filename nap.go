package nap

import (
	"encoding/json"
	"net/http"
)

var (
	MethodNotAllowedHandler = HandlerFunc(func(req *http.Request) (interface{}, Status) {
		return nil, MethodNotAllowed{"method not allowed on resource"}
	})
	NotFoundHandler = HandlerFunc(func(req *http.Request) (interface{}, Status) {
		return nil, NotFound{"resource not found"}
	})
)

type Wrapper interface {
	Wrap(payload json.Marshaler, status Status) json.Marshaler
}

type DefaultWrapper struct{}

func (w DefaultWrapper) Wrap(payload interface{}, status Status) map[string]interface{} {
	return map[string]interface{}{
		"status": map[string]interface{}{
			"code":    status.Code(),
			"message": status.Message(),
		},
		"result": payload,
	}

}

type HandlerFunc func(req *http.Request) (interface{}, Status)

func (f HandlerFunc) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	payload, status := f(request)
	result := DefaultWrapper{}.Wrap(payload, status)
	if res, err := json.Marshal(result); err == nil {
		writer.Header().Add("Content-Type", "application/json")
		writer.Header().Add("Cache-Control", "no-cache,must-revalidate")
		writer.WriteHeader(status.Code())
		writer.Write(res)
	} else {
		writer.WriteHeader(http.StatusInternalServerError)
	}
}
