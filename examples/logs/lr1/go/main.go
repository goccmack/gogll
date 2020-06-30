package main

import (
	"fmt"
	"logs/lexer"
	"logs/parser"
	"time"
)

const inputFile = "../../data/test.log"

func main() {
	start := time.Now()
	lex := lexer.NewFile(inputFile)
	lexDone := time.Now()
	_, err := parser.New(lex).Parse()
	if err != nil {
		panic(err)
	}
	parseDone := time.Now()
	fmt.Printf("Lexer duration %d ms\n", lexDone.Sub(start)/time.Millisecond)
	fmt.Printf("Parser duration %d ms\n", parseDone.Sub(lexDone)/time.Millisecond)
}
