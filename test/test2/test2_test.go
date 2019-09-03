package test2

import (
	"github.com/goccmack/gogll/test/test2/goutil/bsr"
	"github.com/goccmack/gogll/test/test2/parser"
	"testing"
)

func Test1(t *testing.T) {
	if err, errs := parser.Parse([]byte("aa")); err != nil {
		t.Error(err)
		for _, err := range errs {
			t.Log(err)
		}
	}
	bsr.Dump()
}