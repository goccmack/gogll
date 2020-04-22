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
	Package        *Package
	LexRules       []*LexRule
	SyntaxRules    []*SyntaxRule
	Terminals      *stringset.StringSet
	NonTerminals   *stringset.StringSet
	StringLiterals *stringset.StringSet
}

type NT struct {
	tok *token.Token
}

type Package struct {
	tok *token.Token
}

type TokID struct {
	tok *token.Token
}

type Terminal interface {
	isTerminal()
}

/*** Methods ***/

func (g *GoGLL) GetLexRule(id string) *LexRule {
	for _, r := range g.LexRules {
		if r.TokID.Token() == id {
			return r
		}
	}
	return nil
}

func (g *GoGLL) GetSyntaxRule(nt string) *SyntaxRule {
	for _, r := range g.SyntaxRules {
		if r.Head.Token() == nt {
			return r
		}
	}
	return nil
}

func (g *GoGLL) GetSymbols() []string {
	return append(g.Terminals.Elements(), g.NonTerminals.Elements()...)
}

func (g *GoGLL) StartSymbol() string {
	return g.SyntaxRules[0].Head.Token()
}

func (n *NT) String() string {
	return string(n.tok.Literal)
}

func (n *NT) Token() string {
	return string(n.tok.Literal)
}

func (n *NT) Lext() int {
	return n.tok.Lext
}

// ID returns the identifier of n
func (n *NT) ID() string {
	return string(n.tok.Literal)
}

func (p *Package) GetString() string {
	return string(p.tok.Literal[1 : len(p.tok.Literal)-1])
}

func (s *StringLit) String() string {
	return string(s.tok.Literal[1 : len(s.tok.Literal)-1])
}

func (s *StringLit) Token() string {
	return string(s.tok.Literal[1 : len(s.tok.Literal)-1])
}

func (s *StringLit) Lext() int {
	return s.tok.Lext
}

func (t *TokID) String() string {
	return string(t.tok.Literal)
}

func (t *TokID) Token() string {
	return string(t.tok.Literal)
}

func (t *TokID) Lext() int {
	return t.tok.Lext
}
