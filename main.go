//  Copyright 2019 Marius Ackerman
//
//  Licensed under the Apache License, Version 2.0 (the "License");
//  you may not use this file except in compliance with the License.
//  You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
//  Unless required by applicable law or agreed to in writing, software
//  distributed under the License is distributed on an "AS IS" BASIS,
//  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//  See the License for the specific language governing permissions and
//  limitations under the License.

package main

import (
	"fmt"
	"gogll/ast"
	"gogll/cfg"
	"gogll/check"
	genff "gogll/gen/firstfollow"
	"gogll/gen/golang"
	"gogll/gen/slots"
	"gogll/gen/symbols"
	"gogll/goutil/md"
	"gogll/lexer"
	"gogll/parser"
	"os"
	"strings"
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
