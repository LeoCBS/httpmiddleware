//go:build unit
// +build unit

package httpmiddleware_test

import (
	"testing"

	"github.com/LeoCBS/httpmiddleware"
	"github.com/sirupsen/logrus"
)

func TestNewMiddlewareWorks(t *testing.T) {
	l := logrus.New()
	md := httpmiddleware.New(l)
	assertNotNil(t, md)
}

func TestNewMiddlewareValidateURLParams(t *testing.T) {
	l := logrus.New()
	md := httpmiddleware.New(l)
	assertNotNil(t, md)
}

func assertNotNil(t *testing.T, value interface{}) {
	t.Helper()
	if value == nil {
		t.Error("value is nil")
	}
}
