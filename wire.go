package main

import (
	"fmt"
	"net/http"
	"reflect"
	"regexp"
)

type Wire struct {
	prefix      string
	routes      []*route
	middlewares []func(http.Handler) http.Handler
}

type route struct {
	pattern *regexp.Regexp
	execs   []*routeExec
}

type routeExec struct {
	method  string
	handler http.Handler
}

func NewWire() *Wire {
	return &Wire{}
}

func (wi *Wire) Chain(middlewares ...func(http.Handler) http.Handler) {
	for i := len(middlewares) - 1; i >= 0; i-- {
		wi.middlewares = append(wi.middlewares, middlewares[i])
	}
}

func (wi *Wire) registerHandler(path string, method string, h http.Handler) {

	re := &routeExec{method: method, handler: h}

	var pathExists bool
	str := "^%s$"
	if path[len(path)-1:] == "/" {
		str = "^%s"
	}
	pattern := regexp.MustCompile(fmt.Sprintf(str, wi.prefix+path))
	for _, r := range wi.routes {
		if reflect.DeepEqual(r.pattern, pattern) {
			r.execs = append(r.execs, re)
			pathExists = true
		}
	}
	if !pathExists {
		wi.routes = append(wi.routes, &route{pattern, []*routeExec{re}})
	}
}

func handlerWithStatus(status int) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(status)
	})
}

func (wi *Wire) lookupHandler(req *http.Request) http.Handler {
	var h http.Handler

	var pathFound, methodAllowed bool
	for _, r := range wi.routes {
		if r.pattern.MatchString(req.URL.Path) {
			pathFound = true
			for _, r := range r.execs {
				if r.method == req.Method || r.method == "ALL" {
					methodAllowed = true
					h = r.handler
					break
				}
			}
			break
		}
	}
	if !pathFound {
		h = handlerWithStatus(http.StatusNotFound)
	}
	if pathFound && !methodAllowed {
		h = handlerWithStatus(http.StatusMethodNotAllowed)
	}
	return h
}

func (wi *Wire) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h := wi.lookupHandler(r)
	for _, m := range wi.middlewares {
		h = m(h)
	}
	h.ServeHTTP(w, r)
}

func (wi *Wire) SubRouter(path string) *Wire {
	nwi := NewWire()
	nwi.prefix = wi.prefix + path
	wi.registerHandler(path+"/", "ALL", nwi)
	return nwi
}
