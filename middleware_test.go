//go:build unit
// +build unit

package main_test

import (
	"testing"

	"github.com/LeoCBS/httpmiddleware"
	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"
)

func TestNewMiddleware(t *testing.T) {
	l := logrus.New()
	router := httprouter.New()
	m := httpmiddleware.New(l, router)
}
