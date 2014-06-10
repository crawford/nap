package nap

import (
	"encoding/json"
	"net/http"
)

var (
	MethodNotAllowedHandler = NapHandlerFunc(func() (interface{}, Status) {
		return nil, MethodNotAllowed{
			message: "method not allowed on resource",
		}
	})
	NotFoundHandler = NapHandlerFunc(func() (interface{}, Status) {
		return nil, NotFound{
			message: "resource not found",
		}
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

type NapHandlerFunc func() (interface{}, Status)

func (f NapHandlerFunc) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	payload, status := f()
	result := DefaultWrapper{}.Wrap(payload, status)
	if res, err := json.Marshal(result); err == nil {
		writer.WriteHeader(status.Code())
		writer.Write(res)
	} else {
		writer.WriteHeader(http.StatusInternalServerError)
	}
}
