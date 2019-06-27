package main

import (
	"flag"
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
	"path"
)

var (
	bnfFile string
	baseDir string
)

const filePerm os.FileMode = 0731

func main() {
	getParameters()
	lex, err := lexer.NewLexerFile(bnfFile)
	if err != nil {
		fail(fmt.Sprintf("%s", err))
	}
	p := parser.NewParser()
	parse, err := p.Parse(lex)
	if err != nil {
		fmt.Printf("PARSE ERROR: %s\n", err)
		os.Exit(1)
	}
	g := ast.GetAST(parse)
	check.Check(g)
	symbols.Gen(baseDir)
	genff.Gen(baseDir, g)
	slots.Gen(baseDir, filePerm)
	golang.Gen(baseDir, g)
}

func getFileBase() {
	baseDir, _ = path.Split(bnfFile)
}

func getSourceFile() {
	if flag.NArg() < 1 {
		fail("Source file required")
	}
	bnfFile = flag.Arg(0)
}

func fail(msg string) {
	fmt.Printf("ERROR: %s\n", msg)
	usage()
	os.Exit(1)
}

func getParameters() {
	flag.Parse()
	getSourceFile()
	getFileBase()
}

func usage() {
	msg := `use: gogll <bnfFile file>
	<file name> : Name of the BNF file to be processed`
	fmt.Println(msg)
}
