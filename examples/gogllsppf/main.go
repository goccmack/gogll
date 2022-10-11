package main

import (
	"github.com/goccmack/gogll/v3/lexer"
	"github.com/goccmack/gogll/v3/parser"
)

func main() {
	lex := lexer.NewFile("../../gogll.md")
	pf, errs := parser.Parse(lex)
	if errs != nil {
		panic(errs)
	}
	pf.ToSPPF().DotFile("gogllSPPF.dot")
}
