package wire

import (
	"net/http"
)

func (rt *Router) Get(p string, h http.Handler) {
	rt.registerHandler(p, http.MethodGet, h)
}

func (rt *Router) GetF(p string, f func(http.ResponseWriter, *http.Request)) {
	rt.registerHandler(p, http.MethodGet, http.HandlerFunc(f))
}

func (rt *Router) Post(p string, h http.Handler) {
	rt.registerHandler(p, http.MethodPost, h)
}

func (rt *Router) PostF(p string, f func(http.ResponseWriter, *http.Request)) {
	rt.registerHandler(p, http.MethodPost, http.HandlerFunc(f))
}

func (rt *Router) All(p string, h http.Handler) {
	rt.registerHandler(p, "ALL", h)
}

func (rt *Router) AllF(p string, f func(http.ResponseWriter, *http.Request)) {
	rt.registerHandler(p, "ALL", http.HandlerFunc(f))
}
