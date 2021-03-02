package main

import (
    "fmt"
	"context"
	"net/http"
	"reflect"
	"regexp"
)

type Router struct {
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

func NewRouter() *Router {
	return &Router{}
}

func (rt *Router) Chain(middlewares ...func(http.Handler) http.Handler) {
	for i := len(middlewares) - 1; i >= 0; i-- {
		rt.middlewares = append(rt.middlewares, middlewares[i])
	}
}

var nRgxp = regexp.MustCompile(`\((.*?):`)

func (rt *Router) normalizePath(path string) string {
	path = nRgxp.ReplaceAllString(path, "(")
	// if path[len(path)-1:] != "/" {
	// 	path += "$"
	// }
	return path
}

var fRgxp = regexp.MustCompile(`\(.*?\)`)

func (rt *Router) findPathVars(path string) []*pathVar {
	nPath := rt.normalizePath(path)

	nms := nRgxp.FindAllStringSubmatch(path, -1)
	fms := fRgxp.FindAllStringSubmatch(nPath, -1)

	var vars []*pathVar
	for i, fm := range fms {
		vars = append(vars, &pathVar{nms[i][1], regexp.MustCompile(fm[0])})
	}

	return vars
}

func (rt *Router) registerHandler(path string, method string, h http.Handler) {

	re := &routeExec{method: method, handler: h}

	nPath := rt.normalizePath(path)
	pattern := regexp.MustCompile("^" + rt.prefix + nPath)

	pathVars := rt.findPathVars(path)

	var pathExists bool
	for _, r := range rt.routes {
		if reflect.DeepEqual(r.pattern, pattern) {
			r.execs = append(r.execs, re)
			pathExists = true
		}
	}
	if !pathExists {
		rt.routes = append(
			rt.routes,
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
            fmt.Println(path)
		}
		r = r.WithContext(ctx)
		h.ServeHTTP(w, r)
	})
}

func (rt *Router) lookupHandler(req *http.Request) http.Handler {
	var h http.Handler

	var pathFound, methodAllowed bool
	for _, r := range rt.routes {
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

func (rt *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h := rt.lookupHandler(r)
	for _, m := range rt.middlewares {
		h = m(h)
	}
	h.ServeHTTP(w, r)
}

func (rt *Router) SubRouter(path string) *Router {
	nwi := NewRouter()
	nwi.prefix = rt.prefix + rt.normalizePath(path)
	rt.registerHandler(path+"/", "ALL", nwi)
	return nwi
}
