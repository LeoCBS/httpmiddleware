//go:build unit
// +build unit

package httpmiddleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/LeoCBS/httpmiddleware"
	"github.com/LeoCBS/httpmiddleware/test"
	"github.com/sirupsen/logrus"
)

func TestNewMiddlewareWorks(t *testing.T) {
	l := logrus.New()
	md := httpmiddleware.New(l)
	test.AssertNotNil(t, md)
}

func TestNewMiddlewareValidateURLParams(t *testing.T) {
	l := logrus.New()
	md := httpmiddleware.New(l)
	test.AssertNotNil(t, md)

	fnHandlePOST := func(w http.ResponseWriter, r *http.Request, ps httpmiddleware.Params) httpmiddleware.Response {
		//TODO here you add your business logic, call some storage
		//func, etc...
		return httpmiddleware.Response{
			StatusCode: http.StatusOK,
		}
	}
	//register a simple route POST using key/value URL pattern
	md.POST("/name/:name/age/:age", fnHandlePOST)
	assertInvalidRequest(t, md)
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
