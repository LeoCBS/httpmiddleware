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

	res := httptest.NewRecorder()
	md.ServeHTTP(res, req)

	test.AssertEqual(t, http.StatusBadRequest, res.Code)
	expectedResponseBody := `{"error":"your URL must inform name value"}`
	test.AssertBodyContains(t, res.Body, expectedResponseBody)
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

	res := httptest.NewRecorder()
	f.md.ServeHTTP(res, req)

	test.AssertEqual(t, http.StatusOK, res.Code)
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
	f.md.POST("/whatever", fnHandlePOST)

	URL := "/whatever"
	req, err := http.NewRequest("POST", URL, nil)
	test.AssertNoError(t, err)

	recorder := httptest.NewRecorder()
	f.md.ServeHTTP(recorder, req)
	resp := recorder.Result()
	test.AssertEqual(t, resp.StatusCode, http.StatusOK)
	test.AssertEqual(t, resp.Header.Get(headerKey), headerValue)
}

func TestNewMiddlewareWriteCustomBadRequest(t *testing.T) {
	f := setUp(t)

	expectedErr := "your body must be one valid JSON"
	fnHandle := func(w http.ResponseWriter, r *http.Request, ps httpmiddleware.Params) httpmiddleware.Response {
		// just returning a errors.BadRequest that middleware will
		// write status code and body response
		return httpmiddleware.Response{
			Error: errors.NewBadRequest(expectedErr),
		}
	}
	f.md.GET("/whatever", fnHandle)

	URL := "/whatever"
	req, err := http.NewRequest(http.MethodGet, URL, nil)
	test.AssertNoError(t, err)

	recorder := httptest.NewRecorder()
	f.md.ServeHTTP(recorder, req)
	resp := recorder.Result()
	test.AssertEqual(t, resp.StatusCode, http.StatusBadRequest)
	expectedResponseBody := fmt.Sprintf(`{"error":"%s"}`, expectedErr)
	test.AssertBodyContains(t, resp.Body, expectedResponseBody)
}

//TODO add test to all HTTP methods / check method not allowed
//TODO check to call result in all tests
