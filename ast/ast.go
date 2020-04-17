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
	"fmt"
	"os"
	"strings"

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

type NT struct {
	tok     *token.Token
	Literal string
}

type Package struct {
	tok *token.Token
}

type Rule struct {
	Head       *NT
	Alternates []*Alternate
}

type StringLit struct {
	tok *token.Token
}

type Symbol interface {
	isSymbol()
	Pos() token.Pos
	Token() string
	String() string
}

func (*NT) isSymbol() {}

// Terminals
func (*TokID) isSymbol()     {}
func (*StringLit) isSymbol() {}

type TokID struct {
	tok *token.Token
}

type Terminal interface {
	isTerminal()
}

func NewGoGLL(pkg, rules interface{}) (*GoGLL, error) {
	gogll := &GoGLL{
		Package: pkg.(*Package),
		Rules:   rules.([]*Rule),
	}
	gogll.NonTerminals = gogll.nonTerminals()
	gogll.Terminals = gogll.terminals()
	return gogll, nil
}

func NewAlternate(alt interface{}) (*Alternate, error) {
	a := &Alternate{}
	if symbols, ok := alt.([]Symbol); ok {
		a.Symbols = symbols
	}
	return a, nil
}

func NewAlternates(alt interface{}) ([]*Alternate, error) {
	alts := []*Alternate{alt.(*Alternate)}
	return alts, nil
}

func AddAlternate(alts, alt interface{}) ([]*Alternate, error) {
	as := append(alts.([]*Alternate), alt.(*Alternate))
	return as, nil
}

func NewNT(nt interface{}) (*NT, error) {
	n := &NT{
		tok: nt.(*token.Token),
	}
	n.Literal = n.tok.IDValue()
	return n, nil
}

func NewPackage(pkg interface{}) (*Package, error) {
	p := &Package{tok: pkg.(*token.Token)}
	return p, nil
}

func NewRule(nt, alts interface{}) (*Rule, error) {
	r := &Rule{
		Head:       nt.(*NT),
		Alternates: alts.([]*Alternate),
	}
	return r, nil
}

func NewRules(rule interface{}) ([]*Rule, error) {
	rs := []*Rule{rule.(*Rule)}
	return rs, nil
}

func AddRule(rules, rule interface{}) ([]*Rule, error) {
	rs := append(rules.([]*Rule), rule.(*Rule))
	return rs, nil
}

func NewStringLit(s interface{}) (*StringLit, error) {
	str := &StringLit{tok: s.(*token.Token)}
	return str, nil
}

func NewSymbols(sym interface{}) ([]Symbol, error) {
	ss := []Symbol{sym.(Symbol)}
	return ss, nil
}

func AddSymbol(symbols, symbol interface{}) ([]Symbol, error) {
	ss := append(symbols.([]Symbol), symbol.(Symbol))
	return ss, nil
}

func NewTokID(tid interface{}) (*TokID, error) {
	tokId := &TokID{
		tok: tid.(*token.Token),
	}
	return tokId, nil
}

/*** Methods ***/

func (a *Alternate) GetSymbols() []string {
	symbols := make([]string, len(a.Symbols))
	for i, s := range a.Symbols {
		symbols[i] = s.Token()
	}
	return symbols
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
	return append(g.terminals().Elements(), g.nonTerminals().Elements()...)
}

func (g *GoGLL) StartSymbol() string {
	return g.Rules[0].Head.Token()
}

func (g *GoGLL) nonTerminals() *stringset.StringSet {
	nts := stringset.New()
	for _, r := range g.Rules {
		if nts.Contain(r.Head.Token()) {
			fail(fmt.Errorf("Duplicate rule %s", r.Head.Token()), r.Head.Pos())
		} else {
			nts.Add(r.Head.Token())
		}
	}
	return nts
}

func (g *GoGLL) terminals() *stringset.StringSet {
	terminals := stringset.New()
	for _, r := range g.Rules {
		for _, a := range r.Alternates {
			for _, s := range a.Symbols {
				switch t := s.(type) {
				case *TokID, *StringLit:
					terminals.Add(t.Token())
				}
			}
		}
	}
	return terminals
}

func (a *Alternate) Empty() bool {
	return len(a.Symbols) == 0
}

// func (a *Any) Pos() token.Pos {
// 	return a.tok.Pos
// }

// func (a *Any) String() string {
// 	return "any"
// }

// func (a *Any) Token() string {
// 	return "any"
// }

// func (*Not) String() string {
// 	return "not"
// }

// func (*Not) Token() string {
// 	return "not"
// }

// func (n *Not) Pos() token.Pos {
// 	return n.tok.Pos
// }

func (n *NT) String() string {
	return n.Literal
}

func (n *NT) Token() string {
	return n.Literal
}

func (n *NT) Pos() token.Pos {
	return n.tok.Pos
}

func (p *Package) GetString() string {
	return p.tok.StringValue()
}

func (s *StringLit) String() string {
	return s.tok.StringValue()
}

func (s *StringLit) Token() string {
	return s.tok.StringValue()
}

func (s *StringLit) Pos() token.Pos {
	return s.tok.Pos
}

func (t *TokID) String() string {
	return t.tok.IDValue()
}

func (t *TokID) Token() string {
	return t.tok.IDValue()
}

func (t *TokID) Pos() token.Pos {
	return t.tok.Pos
}

/*** Util ***/

// parse not "..." or any "..." and return a set representing the characters
func stringToSet(str string) *stringset.StringSet {
	i := strings.Index(str, "\"")
	rs := []rune(str[i:])
	ss := stringset.New()
	for j := 0; j < len(rs)-1; {
		r := rs[j]
		j++
		if r == '\\' {
			r = rs[j]
			j++
			switch r {
			case '\\':
				ss.Add("\\")
			case 'n':
				ss.Add("\n")
			case 'r':
				ss.Add("\r")
			case 't':
				ss.Add("\t")
			case '"':
				ss.Add("\"")
			default:
				panic(fmt.Sprintf("Invalid %c", r))
			}
		} else {
			ss.Add(fmt.Sprintf("%s", r))
		}
	}
	return ss
}

/*** Errors ***/

func fail(err error, pos token.Pos) {
	fmt.Printf("Error: %s at line %d col %d\n", err, pos.Line, pos.Column)
	os.Exit(1)
}
