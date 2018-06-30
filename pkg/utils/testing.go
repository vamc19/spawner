package utils

import (
	"testing"
)

func AssertEqual(t *testing.T, a interface{}, b interface{}, msg string) {
	if a != b {
		t.Errorf("%v != %v : %s", a, b, msg)
	}
}

func AssertNotEqual(t *testing.T, a interface{}, b interface{}, msg string) {
	if a == b {
		t.Errorf("%v == %v : %s", a, b, msg)
	}
}
