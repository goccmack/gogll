package test2

import (
	"gogll/test/test2/parser"
	"testing"
)

func Test1(t *testing.T) {
	if err := parser.Parse([]byte("a.b")); err != nil {
		t.Error(err)
	}
}