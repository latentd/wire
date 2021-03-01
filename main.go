package main

import (
    "fmt"
    "net/http"
)

func l1(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        fmt.Println("l1")
        next.ServeHTTP(w, r)
    })
}

func l2(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        fmt.Println("l2")
        next.ServeHTTP(w, r)
    })
}

func t() http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("test"))
    })
}

func t2(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("test2"))
}

func main() {

    w := NewWire()

    w.Chain(l1, l2)
    w.Get("/test", t())

    w2 := w.SubRouter("/api")
    w2.GetF("/a", t2)
    w2.GetF("/a/([1-9]+)/b", t2)
    //r2.GetF("/a/(id:[1-9]+)/b/(nid:[1-9]+)/c", t2)

    http.ListenAndServe(":8080", w)
}
