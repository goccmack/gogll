package example01

import (
	"fmt"
	"gamma1/lexer"
	"gamma1/parser"
	"os"
	"testing"
)

func Test1(t *testing.T) {
	pf, errs := parser.Parse(lexer.New([]rune("aab")))
	if errs != nil {
		parseErrors(errs)
		panic(errs)
	}

	pf.Dump()

	if pf.IsAmbiguous() {
		pf.ReportAmbiguous()
	}
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
