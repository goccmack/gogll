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
Package ast is an Abstract Syntax Tree for gogll, used for code generation.
*/
package ast

import (
	"github.com/goccmack/gogll/token"
	"github.com/goccmack/goutil/stringset"
)

type GoGLL struct {
	Package      *Package
	Rules        []*Rule
	Terminals    *stringset.StringSet
	NonTerminals *stringset.StringSet
}

type Alternate struct {
	Symbols []Symbol
}

type Any struct {
	Tok *token.Token
}

type Not struct {
	Tok *token.Token
}

type NT struct {
	Tok *token.Token
}

type Package struct {
	Tok *token.Token
}

type Rule struct {
	Head       *NT
	Alternates []*Alternate
}

type StringLit struct {
	Tok *token.Token
}

type Symbol interface {
	isSymbol()
	// Lext returns the left extent of Symbol in the input string
	Lext() int
	Token() string
	String() string
}

func (*NT) isSymbol() {}

// Terminals
func (*Any) isSymbol()       {}
func (*Not) isSymbol()       {}
func (*TokID) isSymbol()     {}
func (*StringLit) isSymbol() {}

type TokID struct {
	Tok *token.Token
}

type Terminal interface {
	isTerminal()
}

/*** Methods ***/

func (a *Alternate) GetSymbols() []string {
	symbols := make([]string, len(a.Symbols))
	for i, s := range a.Symbols {
		symbols[i] = s.Token()
	}
	return symbols
}

func (a *Any) Lext() int {
	return a.Tok.Lext
}

func (a *Any) String() string {
	return "any"
}

func (a *Any) Token() string {
	return "any"
}

func (n *Not) Lext() int {
	return n.Tok.Lext
}

func (n *Not) String() string {
	return "not"
}

func (n *Not) Token() string {
	return "not"
}

func (g *GoGLL) GetRule(nt string) *Rule {
	for _, r := range g.Rules {
		if r.Head.Token() == nt {
			return r
		}
	}
	panic("No rule " + nt)
}

func (g *GoGLL) GetSymbols() []string {
	return append(g.Terminals.Elements(), g.NonTerminals.Elements()...)
}

func (g *GoGLL) StartSymbol() string {
	return g.Rules[0].Head.Token()
}

func (a *Alternate) Empty() bool {
	return len(a.Symbols) == 0
}

func (n *NT) String() string {
	return string(n.Tok.Literal)
}

func (n *NT) Token() string {
	return string(n.Tok.Literal)
}

func (n *NT) Lext() int {
	return n.Tok.Lext
}

// ID returns the identifier of n
func (n *NT) ID() string {
	return string(n.Tok.Literal)
}

func (p *Package) GetString() string {
	return string(p.Tok.Literal[1 : len(p.Tok.Literal)-1])
}

func (s *StringLit) String() string {
	return string(s.Tok.Literal[1 : len(s.Tok.Literal)-1])
}

func (s *StringLit) Token() string {
	return string(s.Tok.Literal[1 : len(s.Tok.Literal)-1])
}

func (s *StringLit) Lext() int {
	return s.Tok.Lext
}

func (t *TokID) String() string {
	return string(t.Tok.Literal)
}

func (t *TokID) Token() string {
	return string(t.Tok.Literal)
}

func (t *TokID) Lext() int {
	return t.Tok.Lext
}
