package example01

import (
	"fmt"
	"os"
	"testing"

	"github.com/goccmack/gogll/v3/test/bsr/test1/lexer"
	"github.com/goccmack/gogll/v3/test/bsr/test1/parser"
	"github.com/goccmack/gogll/v3/test/bsr/test1/parser/bsr"
)

func Test1(t *testing.T) {
	pf, errs := parser.Parse(lexer.New([]rune("a b c a b c")))
	if errs != nil {
		parseErrors(errs)
	}

	// pf.Dump()

	if pf.IsAmbiguous() {
		pf.ReportAmbiguous()
		t.Fatal()
	}

	testRoot(pf.GetRoot(), t)

	// fmt.Println("SPPF:")
	// spf := pf.ToSPPF()
	// fmt.Println(spf)
}

func parseErrors(errs []*parser.Error) {
	fmt.Println("Parse Errors:")
	ln := errs[0].Line
	for _, err := range errs {
		if err.Line == ln {
			fmt.Println(err)
		}
	}
	os.Exit(1)
}

// S : A B C |  A B C S ;
func testRoot(b bsr.BSR, t *testing.T) {
	if b.Alternate() == 0 {
		t.Fatal()
	}
	testA(b.GetNTChildI(0), t)
	testB(b.GetNTChildI(1), t)
	testC(b.GetNTChildI(2), t)
}

// A : "a" ;
func testA(b bsr.BSR, t *testing.T) {
	tok := b.GetTChildI(0)
	if tok.LiteralString() != "a" {
		t.Fatal()
	}
}

// B : "b" ;
func testB(b bsr.BSR, t *testing.T) {
	tok := b.GetTChildI(0)
	if tok.LiteralString() != "b" {
		t.Fatal()
	}
}

// C : "c" ;
func testC(b bsr.BSR, t *testing.T) {
	tok := b.GetTChildI(0)
	if tok.LiteralString() != "c" {
		t.Fatal()
	}
}
