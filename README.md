#Nap#
[![Build Status](https://travis-ci.org/crawford/nap.svg?branch=master)](https://travis-ci.org/crawford/nap)
[![Coverage Status](https://coveralls.io/repos/crawford/nap/badge.png?branch=master)](https://coveralls.io/r/crawford/nap?branch=master)

Simple, and I mean simple, REST framework for JSON-based responses.

##Usage##

###Simple###

Here is a dead-simple REST API.

```
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
```

And here it is in action!

```
{
	"result": "Hello, World!",
	"status": {
		"code": 200,
		"message": ""
	}
}
```

```
{
	"result": null,
	"status": {
		"code": 404,
		"message": "resource not found"
	}
}
```

###Extra Parameters###

Here is a more interesting example. This pattern is useful for injecting extra parameters into your handlers.

```
package main

import (
        "github.com/crawford/nap"
        "net/http"
        "net/url"
        "time"
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

```

And the result!

```
{
    "result": {
        "Url": {
            "Scheme": "",
            "Opaque": "",
            "User": null,
            "Host": "",
            "Path": "/info",
            "RawQuery": "",
            "Fragment": ""
        },
        "Host": "localhost:8080",
        "Time": "2014-06-21T23:26:09.002024789-07:00"
    },
    "status": {
        "code": 200,
        "message": ""
    }
}
```

###Panic###

Not a very good programmer?

```
package main

import (
        "github.com/crawford/nap"
        "net/http"
)

func Panic(req *http.Request) (interface{}, nap.Status) {
        panic("AHH")
}

func main() {
        http.Handle("/panic", nap.HandlerFunc(Panic))
        http.Handle("/", nap.NotFoundHandler)
        http.ListenAndServe(":8080", nil)
}
```

Nice save!

```
{
    "result": null,
    "status": {
        "code": 500,
        "message": ""
    }
}
```
