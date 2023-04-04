//go:build unit
// +build unit

package httpmiddleware_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/LeoCBS/httpmiddleware"
	"github.com/LeoCBS/httpmiddleware/errors"
	"github.com/LeoCBS/httpmiddleware/test"
	"github.com/sirupsen/logrus"
)

func TestNewMiddlewareWorks(t *testing.T) {
	l := logrus.New()
	md := httpmiddleware.New(l)
	test.AssertNotNil(t, md)
}

type fixture struct {
	md *httpmiddleware.Middleware
}

func setUp(t *testing.T) *fixture {
	l := logrus.New()
	md := httpmiddleware.New(l)
	test.AssertNotNil(t, md)
	return &fixture{
		md: md,
	}
}

func TestMiddlewareValidateURLParams(t *testing.T) {
	f := setUp(t)

	fnHandlePOST := func(w http.ResponseWriter, r *http.Request, ps httpmiddleware.Params) httpmiddleware.Response {
		//TODO here you add your business logic, call some storage
		//func, etc...
		return httpmiddleware.Response{
			StatusCode: http.StatusOK,
		}
	}
	//register a simple route POST using key/value URL pattern
	f.md.POST("/name/:name/age/:age", fnHandlePOST)
	assertInvalidRequest(t, f.md)
}

func assertInvalidRequest(t *testing.T, md *httpmiddleware.Middleware) {
	URLwithoutNameValue := "/name//age/17"
	req, err := http.NewRequest("POST", URLwithoutNameValue, nil)
	test.AssertNoError(t, err)

	recorder := httptest.NewRecorder()
	md.ServeHTTP(recorder, req)
	resp := recorder.Result()

	test.AssertEqual(t, http.StatusBadRequest, resp.StatusCode)
	expectedResponseBody := `{"error":"your URL must inform name value"}`
	test.AssertBodyContains(t, resp.Body, expectedResponseBody)
}

func TestMiddlewareParseURLParameters(t *testing.T) {
	f := setUp(t)

	var receivedNameValue string
	var receivedAgeValue string
	fnHandlePOST := func(w http.ResponseWriter, r *http.Request, ps httpmiddleware.Params) httpmiddleware.Response {
		//TODO here you add your business logic, call some storage
		//func, etc...
		receivedNameValue = ps.ByName("name")
		receivedAgeValue = ps.ByName("age")
		return httpmiddleware.Response{
			StatusCode: http.StatusOK,
		}
	}
	//register a simple route POST using key/value URL pattern
	f.md.POST("/name/:name/age/:age", fnHandlePOST)

	nameParam := "leo"
	ageParam := "17"
	URL := fmt.Sprintf("/name/%s/age/%s", nameParam, ageParam)
	req, err := http.NewRequest("POST", URL, nil)
	test.AssertNoError(t, err)

	recorder := httptest.NewRecorder()
	f.md.ServeHTTP(recorder, req)
	resp := recorder.Result()

	test.AssertEqual(t, http.StatusOK, resp.StatusCode)
	test.AssertEqual(t, receivedNameValue, nameParam)
	test.AssertEqual(t, receivedAgeValue, ageParam)
}

func TestNewMiddlewareWriteResponseHeaders(t *testing.T) {
	f := setUp(t)

	headerKey := "Location"
	headerValue := "/whatever/01234"
	respHeaders := map[string]string{
		headerKey: headerValue,
	}
	fnHandlePOST := func(w http.ResponseWriter, r *http.Request, ps httpmiddleware.Params) httpmiddleware.Response {
		//TODO here you add your business logic, call some storage
		//func, etc...
		return httpmiddleware.Response{
			StatusCode: http.StatusOK,
			Headers:    respHeaders,
		}
	}
	URL := "/responseheaders"
	f.md.POST(URL, fnHandlePOST)

	req, err := http.NewRequest("POST", URL, nil)
	test.AssertNoError(t, err)

	recorder := httptest.NewRecorder()
	f.md.ServeHTTP(recorder, req)
	resp := recorder.Result()
	test.AssertEqual(t, resp.StatusCode, http.StatusOK)
	test.AssertEqual(t, resp.Header.Get(headerKey), headerValue)
}

func TestMiddlewareHandlingClientErrors(t *testing.T) {
	type newErrorFn func(s string) error
	type tcase struct {
		tName              string
		expectedErrMsg     string
		errFunc            newErrorFn
		expectedStatusCode int
	}

	cases := []tcase{
		tcase{
			tName:              "handlingBadRequest",
			expectedErrMsg:     "your body must be one valid JSON",
			expectedStatusCode: http.StatusBadRequest,
			errFunc:            errors.NewBadRequest,
		},
		tcase{
			tName:              "handlingNotFound",
			expectedErrMsg:     "record not found",
			expectedStatusCode: http.StatusNotFound,
			errFunc:            errors.NewNotFound,
		},
	}
	for _, c := range cases {
		t.Run(c.tName, func(t *testing.T) {
			f := setUp(t)
			fnHandle := func(w http.ResponseWriter, r *http.Request, ps httpmiddleware.Params) httpmiddleware.Response {
				// just returning a github.com/LeoCBS/httpmiddleware/errors that middleware will
				// write status code and body response
				return httpmiddleware.Response{
					Error: c.errFunc(c.expectedErrMsg),
				}
			}
			URL := "/clienterrorhandling"
			f.md.GET(URL, fnHandle)

			req, err := http.NewRequest(http.MethodGet, URL, nil)
			test.AssertNoError(t, err)

			recorder := httptest.NewRecorder()
			f.md.ServeHTTP(recorder, req)
			resp := recorder.Result()
			test.AssertEqual(t, resp.StatusCode, c.expectedStatusCode)
			expectedResponseBody := fmt.Sprintf(`{"error":"%s"}`, c.expectedErrMsg)
			test.AssertBodyContains(t, resp.Body, expectedResponseBody)

		})
	}
}

func TestMiddlewareHandlingInternalServerErrors(t *testing.T) {
	type newErrorFn func(s string) error
	type tcase struct {
		tName              string
		logMsg             string
		expectedErrMsg     string
		errFunc            newErrorFn
		expectedStatusCode int
	}

	cases := []tcase{
		tcase{
			tName:              "handlingInternalServerError",
			expectedErrMsg:     "Internal Server Error",
			logMsg:             "msg used just on log",
			expectedStatusCode: http.StatusInternalServerError,
			errFunc:            errors.NewInternalServerError,
		},
		tcase{
			tName:              "handlingGolangError",
			expectedErrMsg:     "Internal Server Error",
			logMsg:             "msg used just on log / error on business logic",
			expectedStatusCode: http.StatusInternalServerError,
			errFunc:            errors.New,
		},
	}
	for _, c := range cases {
		t.Run(c.tName, func(t *testing.T) {
			f := setUp(t)
			fnHandle := func(w http.ResponseWriter, r *http.Request, ps httpmiddleware.Params) httpmiddleware.Response {
				// just returning a github.com/LeoCBS/httpmiddleware/errors that middleware will
				// write status code and body response
				return httpmiddleware.Response{
					Error: c.errFunc(c.logMsg),
				}
			}
			URL := "/internalerrorhandling"
			f.md.GET(URL, fnHandle)

			req, err := http.NewRequest(http.MethodGet, URL, nil)
			test.AssertNoError(t, err)

			recorder := httptest.NewRecorder()
			f.md.ServeHTTP(recorder, req)
			resp := recorder.Result()
			test.AssertEqual(t, resp.StatusCode, c.expectedStatusCode)
			expectedResponseBody := fmt.Sprintf(`{"error":"%s"}`, c.expectedErrMsg)
			test.AssertBodyContains(t, resp.Body, expectedResponseBody)
		})
	}
}

func TestMiddlewareHandlingAllHTTPMethods(t *testing.T) {
	type tcase struct {
		tName              string
		httpMethod         string
		expectedStatusCode int
	}

	cases := []tcase{
		tcase{
			tName:              "handlingGET",
			httpMethod:         http.MethodGet,
			expectedStatusCode: http.StatusOK,
		},
		tcase{
			tName:              "handlingPOST",
			httpMethod:         http.MethodPost,
			expectedStatusCode: http.StatusOK,
		},
		tcase{
			tName:              "handlingPUT",
			httpMethod:         http.MethodPut,
			expectedStatusCode: http.StatusOK,
		},
		tcase{
			tName:              "handlingDELETE",
			httpMethod:         http.MethodGet,
			expectedStatusCode: http.StatusOK,
		},
		tcase{
			tName:              "handlingOPTIONS",
			httpMethod:         http.MethodOptions,
			expectedStatusCode: http.StatusOK,
		},
	}
	for _, c := range cases {
		t.Run(c.tName, func(t *testing.T) {
			f := setUp(t)
			fnHandle := func(w http.ResponseWriter, r *http.Request, ps httpmiddleware.Params) httpmiddleware.Response {
				return httpmiddleware.Response{
					StatusCode: http.StatusOK,
				}
			}
			URL := "/handlingwhatever"
			switch method := c.httpMethod; method {
			case http.MethodGet:
				f.md.GET(URL, fnHandle)
			case http.MethodPost:
				f.md.POST(URL, fnHandle)
			case http.MethodPut:
				f.md.PUT(URL, fnHandle)
			case http.MethodDelete:
				f.md.DELETE(URL, fnHandle)
			case http.MethodOptions:
				f.md.OPTIONS(URL, fnHandle)
			}

			req, err := http.NewRequest(c.httpMethod, URL, nil)
			test.AssertNoError(t, err)

			recorder := httptest.NewRecorder()
			f.md.ServeHTTP(recorder, req)
			resp := recorder.Result()
			test.AssertEqual(t, resp.StatusCode, c.expectedStatusCode)
		})
	}
}

func TestNewMiddlewareWriteResponseHeadersAndResponseBody(t *testing.T) {
	f := setUp(t)

	headerKey := "Location"
	headerValue := "/whatever/01234"
	respHeaders := map[string]string{
		headerKey: headerValue,
	}
	type myBusiness struct {
		Name string `json:"name"`
	}
	expectedBody := myBusiness{Name: "leo"}
	handleFn := func(w http.ResponseWriter, r *http.Request, ps httpmiddleware.Params) httpmiddleware.Response {
		//TODO here you add your business logic, call some storage
		//func, etc...
		return httpmiddleware.Response{
			StatusCode: http.StatusOK,
			Headers:    respHeaders,
			Body:       expectedBody,
		}
	}
	URL := "/success"
	f.md.GET(URL, handleFn)

	req, err := http.NewRequest(http.MethodGet, URL, nil)
	test.AssertNoError(t, err)

	recorder := httptest.NewRecorder()
	f.md.ServeHTTP(recorder, req)
	resp := recorder.Result()
	test.AssertEqual(t, resp.StatusCode, http.StatusOK)
	test.AssertEqual(t, resp.Header.Get(headerKey), headerValue)
	test.AssertBodyContainsStruct(t, resp.Body, expectedBody)
}
