// Package test has functions to help commons asserts
package test

import (
	"fmt"
	"io"
	"reflect"
	"strings"
	"testing"
)

func AssertNotNil(t *testing.T, value interface{}) {
	t.Helper()
	if value == nil {
		t.Error("value is nil")
	}
}

func AssertNoError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Error("expected nil error but got this error = %w", err)
	}
}

func AssertEqual(t *testing.T, obj1 interface{}, obj2 interface{}) {
	t.Helper()
	if !reflect.DeepEqual(obj1, obj2) {
		t.Errorf("obj {%v} not equal to obj {%v}", obj1, obj2)
	}
}

func AssertBodyContains(
	t *testing.T,
	bodyReader io.Reader,
	expectedMsg string,
) {
	t.Helper()
	body, err := io.ReadAll(bodyReader)
	AssertNoError(t, err)
	AssertContains(
		t,
		string(body),
		expectedMsg,
	)
}

func AssertContains(t *testing.T, s, substr string) {
	t.Helper()
	if !strings.Contains(s, substr) {
		detail := fmt.Sprintf("{%s} doesn't contains {%s}", s, substr)
		t.Error(detail)
	}
}
