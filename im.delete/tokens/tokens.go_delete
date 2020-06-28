/*
Copyright 2020 Marius Ackerman

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

  http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Package tokens is the intermediate representation of tokens for all code generators.
package tokens

import (
	"fmt"

	"github.com/goccmack/gogll/ast"
)

type Tokens struct {
	TypeToString    []string
	StringToType    map[string]int
	LiteralToString map[string]string
	TypeToLiteral   []string
}

func New(g *ast.GoGLL) *Tokens {
	tokens := &Tokens{
		TypeToString:    []string{},
		StringToType:    map[string]int{},
		LiteralToString: map[string]string{},
	}
	tokens.add("Error", "Error")
	tokens.add("EOF", "EOF")
	for i, tok := range g.Terminals.ElementsSorted() {
		tokens.add(fmt.Sprintf("Type%d", i), tok)
	}
	return tokens
}

func (t *Tokens) add(tokID, tokLit string) {
	tokType := len(t.TypeToString)
	t.StringToType[tokID] = tokType
	t.TypeToString = append(t.TypeToString, tokID)
	t.TypeToLiteral = append(t.TypeToLiteral, tokLit)
	t.LiteralToString[tokLit] = tokID
}
