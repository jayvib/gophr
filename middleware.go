package main

import (
	"net/http"
	"io"
	"fmt"
	"time"
)

type decoratorHandler func(handler http.Handler) http.Handler
type decoratorHandlerFunc func(handlerFunc http.HandlerFunc) http.HandlerFunc

func NewMiddleware(handler http.Handler) *Middleware {
	return &Middleware{
		handler: handler,
		middlewares: make([]decoratorHandler, 0),
	}
}

type Middleware struct {
	handler http.Handler
	middlewares []decoratorHandler
}

func (m *Middleware) Add(mw ...decoratorHandler) {
	m.middlewares = append(m.middlewares, mw...)
}

func (m *Middleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for _, mw := range m.middlewares {
		m.handler = mw(m.handler)
	}
	m.handler.ServeHTTP(w, r)
}

func LoggingMiddleware(writer io.Writer) func(handler http.Handler) http.Handler {
	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
			fmt.Fprintf(writer, "METHOD: %s | URL: %s | Timestamp: %v\n", r.Method, r.URL.Path, time.Now())
		})
	}
}