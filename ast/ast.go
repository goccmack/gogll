//  Copyright 2020 Marius Ackerman
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
	"sort"

	"github.com/goccmack/gogll/v3/token"
	"github.com/goccmack/goutil/stringset"
)

type GoGLL struct {
	Package        *Package
	LexRules       []*LexRule
	SyntaxRules    []*SyntaxRule
	Terminals      *stringset.StringSet
	NonTerminals   *stringset.StringSet
	StringLiterals map[string]*StringLit
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
		if r.TokID.ID() == id {
			return r
		}
	}
	return nil
}

// GetStringLiterals returns a sorted slice of the string literals
func (g *GoGLL) GetStringLiterals() []string {
	slits := make([]string, 0, len(g.StringLiterals))
	for _, sl := range g.StringLiterals {
		slits = append(slits, string(sl.Value()))
	}
	sort.Slice(slits, func(i, j int) bool { return slits[i] < slits[j] })
	return slits
}

// GetStringLiteralsSet returns a stringset containing the string literals
func (g *GoGLL) GetStringLiteralsSet() *stringset.StringSet {
	return stringset.New(g.GetStringLiterals()...)
}

func (g *GoGLL) GetSyntaxRule(nt string) *SyntaxRule {
	for _, r := range g.SyntaxRules {
		if r.Head.ID() == nt {
			return r
		}
	}
	return nil
}

func (g *GoGLL) GetSymbols() []string {
	return append(g.Terminals.Elements(), g.NonTerminals.Elements()...)
}

func (g *GoGLL) StartSymbol() string {
	return g.SyntaxRules[0].Head.ID()
}

func (n *NT) String() string {
	return n.tok.LiteralString()
}

func (n *NT) Lext() int {
	return n.tok.Lext()
}

// ID returns the identifier of n
func (n *NT) ID() string {
	return n.tok.LiteralString()
}

func (p *Package) GetString() string {
	return string(p.tok.Literal()[1 : len(p.tok.Literal())-1])
}

func (s *StringLit) String() string {
	return string(s.tok.Literal()[1 : len(s.tok.Literal())-1])
}

func (s *StringLit) Lext() int {
	return s.tok.Lext()
}

func (t *TokID) String() string {
	return t.tok.LiteralString()
}

func (t *TokID) Lext() int {
	return t.tok.Lext()
}

func (t *TokID) ID() string {
	return t.tok.LiteralStringStripEscape()
}
