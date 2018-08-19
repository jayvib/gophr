package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
)

type Adapter func(handler http.Handler) http.Handler

type Middleware []http.Handler

func (m *Middleware) Add(handler http.Handler) {
	*m = append(*m, handler)
}

func (m *Middleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	mw := NewMiddlewareResponseWriter(w)
	for _, handler := range *m {
		fmt.Println("Running handler")
		handler.ServeHTTP(mw, r)
		if mw.written {
			fmt.Println("written")
			return
		}
	}
	fmt.Fprintln(w, "some thing bad happen")
}

type MiddlewareResponseWriter struct {
	http.ResponseWriter
	written bool
}

func (w *MiddlewareResponseWriter) Write(bytes []byte) (int, error) {
	w.written = true
	n, err := w.ResponseWriter.Write(bytes)
	return n, err
}

func (w *MiddlewareResponseWriter) WriteHeader(code int) {
	w.written = true
	w.ResponseWriter.WriteHeader(code)
}

func NewMiddlewareResponseWriter(w http.ResponseWriter) *MiddlewareResponseWriter {
	return &MiddlewareResponseWriter{
		ResponseWriter: w,
	}
}

func AuthMiddleware(h http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if RequestUser(r) != nil {
			fmt.Println("has valid session")
			h.ServeHTTP(w, r)
			return
		}
		fmt.Println("need to login!")
		query := url.Values{}
		query.Add("next", url.QueryEscape(r.URL.String()))
		http.Redirect(w, r, "/login?"+query.Encode(), http.StatusFound)
	})
}

func AuthMiddleware2(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	if RequestUser(r) != nil {
		fmt.Println("has valid session")
		next.ServeHTTP(w, r)
		return
	}
	fmt.Println("need to login!")
	query := url.Values{}
	query.Add("next", url.QueryEscape(r.URL.String()))
	http.Redirect(w, r, "/login?"+query.Encode(), http.StatusFound)
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Do stuff here
		log.Println(r.RequestURI)
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}

func Adapt(h http.Handler, adapters ...Adapter) http.Handler {
	for _, a := range adapters {
		h = a(h)
	}
	return h
}
