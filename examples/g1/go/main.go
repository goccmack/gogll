package main

import (
	"fmt"
	"g1/lexer"
	"g1/parser"
)

func main() {
	lex := lexer.New([]rune("a | b & c | d"))
	_, errs := parser.Parse(lex)
	if len(errs) > 0 {
		fail(errs)
	}
}

func fail(errs []*parser.Error) {
	fmt.Println("Parse Errors:")
	ln := errs[0].Line
	for _, err := range errs {
		if err.Line == ln {
			fmt.Println("  ", err)
		}
	}
}
