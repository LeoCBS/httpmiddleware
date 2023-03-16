package httpmiddleware

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type Log interface {
	Info(args ...interface{})
	Debug(args ...interface{})
	Warn(args ...interface{})
	Error(args ...interface{})
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
