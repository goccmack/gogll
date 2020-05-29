package main

import (
	"fmt"
	"g1/lexer"
	"g1/parser"
	"os"
	"time"
)

func main() {
	// lex := lexer.New([]rune("a | b & c | d & e | f & g | h & i | j"))
	lex := lexer.New([]rune(
		`a1 & a2 | a3 & a4 | a5 & a6 | a7 & a8 | a9 & a10 | a11 & a12 | 
        a13 & a14 | a15 & a16`))
	start := time.Now()
	bsrSet, errs := parser.Parse(lex)
	if len(errs) > 0 {
		fail(errs)
	}
	fmt.Printf("%d μs\n", time.Now().Sub(start)/time.Microsecond)
	fmt.Printf("%d BSRs\n", len(bsrSet.GetAll()))
}

func fail(errs []*parser.Error) {
	fmt.Println("Parse Errors:")
	ln := errs[0].Line
	for _, err := range errs {
		if err.Line == ln {
			fmt.Println("  ", err)
		}
	}
	os.Exit(1)
}
