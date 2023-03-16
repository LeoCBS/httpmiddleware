//go:build unit
// +build unit

package httpmiddleware_test

import (
	"testing"

	"github.com/LeoCBS/httpmiddleware"
	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"
)

func TestNewMiddleware(t *testing.T) {
	l := logrus.New()
	router := httprouter.New()
	_ = httpmiddleware.New(l, router)
}
