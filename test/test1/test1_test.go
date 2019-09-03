package test1

import (
	"testing"
	"github.com/goccmack/gogll/test/test1/parser"
	"github.com/goccmack/gogll/test/test1/goutil/bsr"
)

func Test1(t *testing.T) {
	if errs := parser.Parse([]byte("a")); errs != nil {
		for _, err := range errs {
			t.Log(err)
			t.Fail()
		}
	}
	bsr.Dump()
}