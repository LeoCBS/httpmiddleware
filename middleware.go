package httpmiddleware

import (
	"encoding/json"
	"fmt"
	"net/http"

	errors "github.com/LeoCBS/httpmiddleware/errors"
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
	Error      error
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

func (ps Params) ByName(name string) string {
	for _, p := range ps {
		if p.Key == name {
			return p.Value
		}
	}
	return ""
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

func (m *Middleware) GET(path string, handler routerHandlerFunc) {
	m.router.GET(path, m.handle(handler))
}

func (m *Middleware) PUT(path string, handler routerHandlerFunc) {
	m.router.PUT(path, m.handle(handler))
}

func (m *Middleware) DELETE(path string, handler routerHandlerFunc) {
	m.router.DELETE(path, m.handle(handler))
}

func (m *Middleware) OPTIONS(path string, handler routerHandlerFunc) {
	m.router.OPTIONS(path, m.handle(handler))
}

func (m *Middleware) handle(next routerHandlerFunc) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		if badRequestResp := isInvalidURLParams(ps); badRequestResp.StatusCode != 0 {
			m.writeResponse(w, badRequestResp)
			return
		}
		resp := next(w, r, convertParams(ps))
		for k, v := range resp.Headers {
			w.Header().Set(k, v)
		}
		if resp.Error != nil {
			switch e := resp.Error.(type) {
			case errors.BadRequest:
				m.writeClientResponse(e, w, resp, http.StatusBadRequest)
			case errors.NotFound:
				m.writeClientResponse(e, w, resp, http.StatusNotFound)
			default:
				// Any error types we don't specifically look out for default
				// to serving a HTTP 500
				m.l.Warn("unexpected error on handle request / error = %v", e)
				resp.Body = getInternalServerErrorBody()
				resp.StatusCode = http.StatusInternalServerError
				m.writeResponse(w, resp)
			}
			return
		}
		m.writeResponse(w, resp)
	}
}

func appendContentTypeJSON(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
}

func convertParams(customParams httprouter.Params) Params {
	params := []Param{}
	for _, v := range customParams {
		params = append(params, Param{
			Key:   v.Key,
			Value: v.Value,
		})
	}
	return params

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
	appendContentTypeJSON(w)
	w.WriteHeader(resp.StatusCode)
	if resp.Body != nil {
		err := json.NewEncoder(w).Encode(resp.Body)
		if err != nil {
			m.l.Warn("error to encode msg / %v", err)
		}
	}
}

func getInternalServerErrorBody() map[string]interface{} {
	return map[string]interface{}{
		"error": http.StatusText(http.StatusInternalServerError),
	}
}

func getClientErrorBody(errStr string) map[string]interface{} {
	return map[string]interface{}{
		"error": errStr,
	}
}

func (m *Middleware) writeClientResponse(
	e error,
	w http.ResponseWriter,
	resp Response,
	statusCode int,
) {
	appendContentTypeJSON(w)
	m.l.Warn("client error / error = %v", e)
	resp.Body = getClientErrorBody(e.Error())
	resp.StatusCode = statusCode
	m.writeResponse(w, resp)
}
