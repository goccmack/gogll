package g2

import (
	"github.com/goccmack/gogll/test/g2/goutil/bsr"
	"github.com/goccmack/gogll/test/g2/parser"
	"testing"
)

func Test1(t *testing.T) {
	err := parser.Parse([]byte("abaa"))
	t.Log(parser.PoppedString())
	t.Log(parser.CRFString())
	bsr.Dump()
	if err != nil {
		t.Error(err)
	}
}