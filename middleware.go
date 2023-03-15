package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type Log interface {
	Info(message string, args ...interface{})
	Debug(message string, args ...interface{})
	Warn(message string, args ...interface{})
	Error(message string, args ...interface{})
}

type Middleware struct {
	l      Log
	router *httprouter.Router
}

// ServeHTTP using httprouter implementation instead default golang
// implementation
func (m *Middleware) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	m.router.ServeHTTP(w, req)
}

func New(l Log, r *httprouter.Router) *Middleware {
	return &Middleware{
		l:      l,
		router: r,
	}
}
