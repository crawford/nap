package main

import (
	"net/http"

	"github.com/crawford/nap"
)

func Panic(req *http.Request) (interface{}, nap.Status) {
	panic("AHH")
}

func main() {
	http.Handle("/panic", nap.HandlerFunc(Panic))
	http.Handle("/", nap.NotFoundHandler)
	http.ListenAndServe(":8080", nil)
}
