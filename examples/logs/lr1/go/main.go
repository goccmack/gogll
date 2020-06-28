package main

import (
	"logs/lexer"
	"logs/parser"
)

const inputFile = "../../../../../logs/data/allvhosts.log"

func main() {
	lex := lexer.NewFile(inputFile)
	_, err := parser.New(lex).Parse()
	if err != nil {
		panic(err)
	}
}
