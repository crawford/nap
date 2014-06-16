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
	PayloadWrapper Wrapper = DefaultWrapper{}
)

type Wrapper interface {
	Wrap(payload interface{}, status Status) interface{}
}

type DefaultWrapper struct{}

func (w DefaultWrapper) Wrap(payload interface{}, status Status) interface{} {
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
	var payload interface{}
	var status Status
	defer func() {
		var result interface{}
		defer func() {
			writer.Header().Add("Content-Type", "application/json")
			writer.Header().Add("Cache-Control", "no-cache,must-revalidate")

			if r := recover(); r == nil {
				if res, err := json.Marshal(result); err == nil {
					writer.WriteHeader(status.Code())
					writer.Write(res)
					return
				}
			}
			writer.WriteHeader(http.StatusInternalServerError)
		}()
		if r := recover(); r != nil {
			status = InternalError{}
		}
		result = PayloadWrapper.Wrap(payload, status)
	}()
	payload, status = f(request)
}
