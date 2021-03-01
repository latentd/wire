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
	ctx := r.Context()

	v := ctx.Value("id")
	s, ok := v.(string)
	fmt.Println(ok)

	v2 := ctx.Value("nid")
	s2, ok := v2.(string)
	fmt.Println(ok)

	w.Write([]byte(s + s2))
}

func main() {

	w := NewWire()

	w.Chain(l1, l2)
	w.Get("/test", t())

	w2 := w.SubRouter("/api/(id:[1-9]+)")

	w3 := w2.SubRouter("/action")
	w3.GetF("/(nid:[1-9]+)", t2)
	//r2.GetF("/a/(id:[1-9]+)/b/(nid:[1-9]+)/c", t2)

	http.ListenAndServe(":8080", w)
}
