package main

import (
	"net/http"
	"net/url"
	"time"

	"github.com/crawford/nap"
)

func Info(now time.Time, req *http.Request) (interface{}, nap.Status) {
	return struct {
		Url  *url.URL
		Host string
		Time time.Time
	}{
		Url:  req.URL,
		Host: req.Host,
		Time: now,
	}, nap.OK{}
}

func AddTimestamp(fn func(time.Time, *http.Request) (interface{}, nap.Status)) http.Handler {
	return nap.HandlerFunc(func(req *http.Request) (interface{}, nap.Status) {
		return fn(time.Now(), req)
	})
}

func main() {
	http.Handle("/info", AddTimestamp(Info))
	http.Handle("/", nap.NotFoundHandler)
	http.ListenAndServe(":8080", nil)
}
