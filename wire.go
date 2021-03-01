package main

import (
	"context"
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
	vars    []*pathVar
	execs   []*routeExec
}

type pathVar struct {
	name    string
	pattern *regexp.Regexp
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

var nRgxp = regexp.MustCompile(`\((.*?):`)

func (wi *Wire) normalizePath(path string) string {
	path = nRgxp.ReplaceAllString(path, "(")
	// if path[len(path)-1:] != "/" {
	// 	path += "$"
	// }
	return path
}

var fRgxp = regexp.MustCompile(`\(.*?\)`)

func (wi *Wire) findPathVars(path string) []*pathVar {
	nPath := wi.normalizePath(path)

	nms := nRgxp.FindAllStringSubmatch(path, -1)
	fms := fRgxp.FindAllStringSubmatch(nPath, -1)

	var vars []*pathVar
	for i, fm := range fms {
		vars = append(vars, &pathVar{nms[i][1], regexp.MustCompile(fm[0])})
	}

	return vars
}

func (wi *Wire) registerHandler(path string, method string, h http.Handler) {

	re := &routeExec{method: method, handler: h}

	nPath := wi.normalizePath(path)
	pattern := regexp.MustCompile("^" + wi.prefix + nPath)

	pathVars := wi.findPathVars(path)

	var pathExists bool
	for _, r := range wi.routes {
		if reflect.DeepEqual(r.pattern, pattern) {
			r.execs = append(r.execs, re)
			pathExists = true
		}
	}
	if !pathExists {
		wi.routes = append(
			wi.routes,
			&route{
				pattern,
				pathVars,
				[]*routeExec{re},
			},
		)
	}
}

func handlerWithStatus(status int) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(status)
	})
}

var aNorm = regexp.MustCompile(`\((.*?)\)`)

func wrapHandlerWithVars(h http.Handler, vars []*pathVar, path string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		for _, v := range vars {
			loc := v.pattern.FindStringIndex(path)
			ctx = context.WithValue(ctx, v.name, path[loc[0]:loc[1]])
			path = path[loc[1]:]
		}
		r = r.WithContext(ctx)
		h.ServeHTTP(w, r)
	})
}

func (wi *Wire) lookupHandler(req *http.Request) http.Handler {
	var h http.Handler

	var pathFound, methodAllowed bool
	for _, r := range wi.routes {
		if r.pattern.MatchString(req.URL.Path) {
			pathFound = true
			for _, e := range r.execs {
				if e.method == req.Method || e.method == "ALL" {
					methodAllowed = true
					h = wrapHandlerWithVars(e.handler, r.vars, req.URL.Path)
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
	nwi.prefix = wi.prefix + wi.normalizePath(path)
	wi.registerHandler(path+"/", "ALL", nwi)
	return nwi
}
