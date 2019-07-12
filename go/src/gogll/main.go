package main

import (
	"gogll/cfg"
	"fmt"
	"gogll/ast"
	"gogll/check"
	genff "gogll/gen/firstfollow"
	"gogll/gen/golang"
	"gogll/gen/slots"
	"gogll/gen/symbols"
	"gogll/lexer"
	"gogll/parser"
	"os"
)

var (
	bnfFile string
	baseDir string
)

func main() {
	cfg.GetParams()
	lex, err := lexer.NewLexerFile(cfg.BNFFile)
	if err != nil {
		fmt.Printf("Error creating lexer: %s", err)
		os.Exit(1)
	}
	p := parser.NewParser()
	parse, err := p.Parse(lex)
	if err != nil {
		fmt.Printf("PARSE ERROR: %s\n", err)
		os.Exit(1)
	}
	g := ast.GetAST(parse)
	check.Check(g)
	symbols.Gen()
	genff.Gen()
	slots.Gen()
	golang.Gen()
}

