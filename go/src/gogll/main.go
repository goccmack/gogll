package main

import (
	"strings"
	"gogll/goutil/md"
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

func main() {
	cfg.GetParams()
	lex := getLexer()
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

func getLexer() *lexer.Lexer {
	if strings.HasSuffix(cfg.SrcFile, ".md") {
		input, err := md.GetSource(cfg.SrcFile)
		if err != nil {
			fmt.Printf("Error extracting source from markdown file: %s", err)
			os.Exit(1)
		}
		return lexer.NewLexer([]byte(input))
	}
	lex, err := lexer.NewLexerFile(cfg.SrcFile)
	if err != nil {
		fmt.Printf("Error creating lexer: %s", err)
		os.Exit(1)
	}
	return lex
}