package empty

import (
	"empty/lexer"
	"empty/parser"
	"strings"
	"testing"
)

const (
	input1 = "aname 123"
	input2 = "123"
)

func Test1(t *testing.T) {
	bs, errs := parser.Parse(lexer.New([]rune(input1)))
	if len(errs) != 0 {
		t.Fail()
	}

	// A1 : Name int ;
	root := bs.GetRoot()

	// Name : name | empty ;
	nm := root.GetNTChildI(0)
	if nm.Alternate() != 0 {
		t.Fail()
	}
	namet := nm.GetTChildI(0)
	if !strings.HasPrefix(input1, namet.LiteralString()) {
		t.Fail()
	}
	numt := root.GetTChildI(1)
	if !strings.HasSuffix(input1, numt.LiteralString()) {
		t.Fail()
	}
}

func Test2(t *testing.T) {
	bs, errs := parser.Parse(lexer.New([]rune(input2)))
	if len(errs) != 0 {
		t.Fail()
	}

	// A1 : Name int ;
	root := bs.GetRoot()

	// Name : name | empty ;
	nm := root.GetNTChildI(0)
	if nm.Alternate() != 1 {
		t.Fail()
	}
	numt := root.GetTChildI(1)
	if input2 != numt.LiteralString() {
		t.Fail()
	}
}
