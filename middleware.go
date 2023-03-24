package httpmiddleware

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type Log interface {
	Info(args ...interface{})
	Debug(args ...interface{})
	Warn(args ...interface{})
	Error(args ...interface{})
}

type Response struct {
	Body       interface{}
	StatusCode int
	Headers    map[string]string
}

type Param struct {
	Key   string
	Value string
}

type Params []Param

type routerHandlerFunc func(w http.ResponseWriter, r *http.Request, ps Params) Response

type Middleware struct {
	l      Log
	router *httprouter.Router
}

// ServeHTTP using httprouter implementation instead default golang
// implementation
func (m *Middleware) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	m.router.ServeHTTP(w, req)
}

func New(l Log) *Middleware {
	router := httprouter.New()
	return &Middleware{
		l:      l,
		router: router,
	}
}

func (m *Middleware) POST(path string, handler routerHandlerFunc) {
	m.router.POST(path, m.handle(handler))
}

func (m *Middleware) handle(next routerHandlerFunc) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		if badRequestResp := isInvalidURLParams(ps); badRequestResp.StatusCode != 0 {
			m.writeResponse(w, badRequestResp)
			return
		}
	}
}

func isInvalidURLParams(ps httprouter.Params) Response {
	for _, p := range ps {
		if p.Value == "" {
			return Response{
				StatusCode: http.StatusBadRequest,
				Body: map[string]interface{}{
					"error": fmt.Sprintf("your URL must inform %s value", p.Key),
				},
			}
		}
	}
	return Response{}
}

func (m *Middleware) writeResponse(w http.ResponseWriter, resp Response) {
	w.WriteHeader(resp.StatusCode)
	if resp.Body != nil {
		err := json.NewEncoder(w).Encode(resp.Body)
		if err != nil {
			m.l.Warn("error to encode msg / %v", err)
		}
	}
}
