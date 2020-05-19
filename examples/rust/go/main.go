package main

import (
	"fmt"
	"gogll/lexer"
	"gogll/parser"
	"time"
)

const N = 10000

func main() {
	start := time.Now()
	var lex *lexer.Lexer
	for i := 0; i < N; i++ {
		lex = lexer.NewFile("../gogll.md")
	}
	lexDone := time.Now()
	for i := 0; i < N; i++ {
		_, errs := parser.Parse(lex)
		if len(errs) > 0 {
			panic("Errors")
		}
	}
	parseDone := time.Now()
	fmt.Printf("Lexer took %d μs\n", lexDone.Sub(start)/(N*time.Microsecond))
	fmt.Printf("Parser took %d μs\n", parseDone.Sub(lexDone)/(N*time.Microsecond))
}
