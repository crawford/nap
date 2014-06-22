package main

import (
	"github.com/crawford/nap"
	"net/http"
)

func Hello(req *http.Request) (interface{}, nap.Status) {
	return "Hello, World!", nap.OK{}
}

func main() {
	http.Handle("/hello", nap.HandlerFunc(Hello))
	http.Handle("/", nap.NotFoundHandler)
	http.ListenAndServe(":8080", nil)
}
