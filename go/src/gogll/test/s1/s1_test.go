package s1

import (
	"testing"
	"gogll/test/s1/parser"
)

func Test1(t *testing.T) {
	if err := parser.Parse([]byte("d")); err != nil {
		t.Error(err)
	}
}

func Test2(t *testing.T) {
	if err := parser.Parse([]byte("daaa")); err != nil {
		t.Error(err)
	}
}