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

/*
Package symbols is a static reposistory of the symbols of generated parser.
It is used by all code generation modules
*/
package symbols

import (
	"fmt"

	"github.com/goccmack/gogll/ast"
)

// Symbol is T or NT
type Symbol interface {
	isSymbol()
	// IsNonTerminal returns true iff this symbol is a non-terminal
	IsNonTerminal() bool

	// Literal returns the literal value of the symbol in the grammar
	Literal() string

	// String returns the string representation of the symbol that is used
	// for code generation
	String() string
}

func (NT) isSymbol() {}
func (T) isSymbol()  {}

// NT is the type of a non-terminal symbol
type NT int

// T is the type of a terminal symbol
type T int

// Symbols is a list of Symbol
type Symbols []Symbol

var (
	initialisized = false
	literalToNT   map[string]NT
	ntToLiteral   []string
	ntToString    []string
	literalToT    map[string]T
	tToLiteral    []string
	tToString     []string
)

// Init initialises the symbols
func Init(g *ast.GoGLL) {
	nts := g.NonTerminals.ElementsSorted()
	ntToLiteral = make([]string, len(nts))
	literalToNT = make(map[string]NT, len(nts))
	ntToString = make([]string, len(nts))
	for i, nt := range nts {
		ntToLiteral[i] = nt
		literalToNT[nt] = NT(i)
		ntToString[i] = fmt.Sprintf("NT_%s", nt)
	}

	ts := g.Terminals.ElementsSorted()
	tToLiteral = make([]string, len(ts))
	literalToT = make(map[string]T, len(ts))
	tToString = make([]string, len(ts))
	for i, t := range ts {
		tToLiteral[i] = t
		literalToT[t] = T(i)
		tToString[i] = fmt.Sprintf("T_%d", i)
	}

	initialisized = true
}

// FromASTString translates an AST symbol string to Symbol
func FromASTString(astSym string) Symbol {
	if !initialisized {
		panic("Uninitialised")
	}
	if nt, exist := literalToNT[astSym]; exist {
		return nt
	}
	if t, exist := literalToT[astSym]; exist {
		return t
	}
	panic(fmt.Sprintf("No symbol %s", astSym))
}

// IsNonTerminal always returns true if the symbol is a non-terminal
func (NT) IsNonTerminal() bool {
	if !initialisized {
		panic("Uninitialised")
	}
	return true
}

// Literal returns the literal value of nt in the grammar
func (nt NT) Literal() string {
	if !initialisized {
		panic("Uninitialised")
	}
	return ntToLiteral[nt]
}

// String returns the string representation of nt used by code generation modules
func (nt NT) String() string {
	if !initialisized {
		panic("Uninitialised")
	}
	return ntToString[nt]
}

// IsNonTerminal always returns false if the symbol is a terminal
func (T) IsNonTerminal() bool {
	if !initialisized {
		panic("Uninitialised")
	}
	return false
}

// Literal returns the literal value of t in the grammar
func (t T) Literal() string {
	if !initialisized {
		panic("Uninitialised")
	}
	return tToLiteral[t]
}

// String returns the string representation of t used by code generation modules
func (t T) String() string {
	if !initialisized {
		panic("Uninitialised")
	}
	return tToString[t]
}

// GetNonTerminals returns the list of non-terminals used by code generation modules
func GetNonTerminals() []NT {
	nts := make([]NT, len(ntToLiteral))
	for i := range nts {
		nts[i] = NT(i)
	}
	return nts
}

// GetTerminals returns the list of terminals used by code generation modules
func GetTerminals() []T {
	ts := make([]T, len(tToLiteral))
	for i := range ts {
		ts[i] = T(i)
	}
	return ts
}

// Empty returns true iff ss is empty
func (ss Symbols) Empty() bool {
	return len(ss) == 0
}

/*
Strings returns a slice containing the string representation of the the
symbols in ss used by code generation modules
*/
func (ss Symbols) Strings() []string {
	strs := make([]string, len(ss))
	for i, s := range ss {
		strs[i] = s.String()
	}
	return strs
}
