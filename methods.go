package main

import (
    "net/http"
)

func (wi *Wire) Get(p string, h http.Handler) {
    wi.registerHandler(p, http.MethodGet, h)
}

func (wi *Wire) GetF(p string, f func(http.ResponseWriter, *http.Request)) {
    wi.registerHandler(p, http.MethodGet, http.HandlerFunc(f))
}

func (wi *Wire) Post(p string, h http.Handler) {
    wi.registerHandler(p, http.MethodPost, h)
}

func (wi *Wire) PostF(p string, f func(http.ResponseWriter, *http.Request)) {
    wi.registerHandler(p, http.MethodPost, http.HandlerFunc(f))
}

func (wi *Wire) All(p string, h http.Handler) {
    wi.registerHandler(p, "ALL", h)
}

func (wi *Wire) AllF(p string, f func(http.ResponseWriter, *http.Request)) {
    wi.registerHandler(p, "ALL", http.HandlerFunc(f))
}
