package httpmiddleware

import (
	"encoding/json"
	"fmt"
	"net/http"

	clienterror "github.com/LeoCBS/httpmiddleware/errors"
	"github.com/julienschmidt/httprouter"
)

type Log interface {
	Info(args ...interface{})
	Debug(args ...interface{})
	Warn(args ...interface{})
	Error(args ...interface{})
}

type Router interface {
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

func New(l Log) *Middleware {
	router := httprouter.New()
	return &Middleware{
		l:      l,
		router: router,
	}
}

type Response struct {
	body       interface{}
	statusCode int
	err        error
	headers    map[string]string
}

type Params httprouter.Params

type routerHandlerFunc func(w http.ResponseWriter, r *http.Request, ps Params) Response

func (m *Middleware) POST(path string, handler routerHandlerFunc) {
	m.router.POST(path, m.handle(handler))
}

func (m *Middleware) handle(next routerHandlerFunc) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		httpRouterParams := convertParams(ps)
		if badRequestResp := isInvalidURLParams(httpRouterParams); badRequestResp.statusCode != 0 {
			m.writeResponse(w, badRequestResp)
			return
		}
		resp := next(w, r, httpRouterParams)
		for k, v := range resp.headers {
			w.Header().Set(k, v)
		}
		if resp.err != nil {
			switch e := resp.err.(type) {
			case clienterror.BadRequestError:
				m.writeClientResponse(e, w, resp, http.StatusBadRequest)
			case clienterror.NotFoundError:
				m.writeClientResponse(e, w, resp, http.StatusNotFound)
			default:
				// Any error types we don't specifically look out for default
				// to serving a HTTP 500
				m.l.Warn("unexpected error on handle request / error = %v", e)
				resp.body = getInternalServerErrorBody()
				m.writeResponse(w, resp)
			}
			return
		}
		m.writeResponse(w, resp)

	}
}

func convertParams(customParams Params) httprouter.Params {
	params := httprouter.Params{}
	for k, v := range customParams {
		params := append(params, httprouter.Param{
			Key:   key,
			Value: v,
		})
	}
	return params

}

func isInvalidURLParams(ps Params) Response {
	for _, p := range ps {
		if p.Value == "" {
			return Response{
				statusCode: http.StatusBadRequest,
				body: map[string]interface{}{
					"error": fmt.Sprintf("your URL must inform %s value", p.Key),
				},
			}
		}
	}
	return Response{}
}

func (m *Middleware) writeResponse(w http.ResponseWriter, resp Response) {
	w.WriteHeader(resp.statusCode)
	if resp.body != nil {
		err := json.NewEncoder(w).Encode(resp.body)
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
	m.l.Warn("client error / error = %v", e)
	resp.body = getClientErrorBody(e.Error())
	resp.statusCode = statusCode
	m.writeResponse(w, resp)
}
