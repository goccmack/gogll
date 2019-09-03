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
	"sort"
	"strings"

	"github.com/goccmack/gogll/goutil/bsr"
	"github.com/goccmack/gogll/goutil/stringset"
)

type Grammar struct {
	Package      string
	Rules        []*Rule
	StartSymbol  string
	CharLiterals []string
	Terminals    map[string]bool
}

func New() *Grammar {
	return &Grammar{
		Terminals: map[string]bool{
			// "any":     true,
			// "letter":  true,
			// "number":  true,
			// "space":   true,
			// "upcase":  true,
			// "lowcase": true,
		},
	}
}

func (g *Grammar) AddPackage(p string) {
	g.Package = p[1 : len(p)-1]
}

func (g *Grammar) AddRule(r *Rule) error {
	if r1 := g.GetRule(r.Head.NT); r1 != nil {
		r1.Append(r)
	} else {
		g.Rules = append(g.Rules, r)
	}
	return nil
}

func (g *Grammar) AddTerminal(t Terminal) {
	g.Terminals[t.String()] = true
}

func (g *Grammar) GetNonTerminals() (nts []string) {
	for _, r := range g.Rules {
		nts = append(nts, r.Head.NT)
	}
	sort.Strings(nts)
	return
}

func (g *Grammar) GetRule(head string) *Rule {
	for _, r := range g.Rules {
		if r.Head.NT == head {
			return r
		}
	}
	return nil
}

func (g *Grammar) SetStartSymbol(nt string) error {
	if g.StartSymbol != "" {
		return fmt.Errorf("Duplicate start symbols %s and %s", g.StartSymbol, nt)
	}
	g.StartSymbol = nt
	return nil
}

func (g *Grammar) GetSymbols() []string {
	return append(g.GetTerminals(), g.GetNonTerminals()...)
}

type Rule struct {
	Head       *Head
	Alternates []*Alternate
}

func (r *Rule) Append(r1 *Rule) error {
	for _, a := range r1.Alternates {
		if r.HasAlternate(a) {
			return fmt.Errorf("Duplicate alternate %s : %s", r.Head, a)
		}
		r.Alternates = append(r.Alternates, a)
	}
	return nil
}

func (r *Rule) GetCharLiterals() (cls []*CharLiteral) {
	for _, a := range r.Alternates {
		cls = append(cls, a.GetCharLiterals()...)
	}
	return
}

func (r *Rule) HasAlternate(a *Alternate) bool {
	for _, a1 := range r.Alternates {
		if a1.Equal(a) {
			return true
		}
	}
	return false
}

type Head struct {
	BSR bsr.BSR
	NT  string
}

func NewHead(b bsr.BSR, nt string) *Head {
	return &Head{
		BSR: b,
		NT:  nt,
	}
}

func (h *Head) String() string {
	return h.NT
}

type Alternate struct {
	Symbols []Symbol
}

func (a *Alternate) Equal(a1 *Alternate) bool {
	if len(a.Symbols) != len(a1.Symbols) {
		return false
	}
	if a.Empty() != a1.Empty() {
		return false
	}
	for i, s := range a.Symbols {
		if !s.Equal(a1.Symbols[i]) {
			return false
		}
	}
	return true
}

func (a *Alternate) GetSymbols() (ss []string) {
	for _, s := range a.Symbols {
		ss = append(ss, s.GetSymbol())
	}
	return
}

func (a *Alternate) GetCharLiterals() (cls []*CharLiteral) {
	for _, s := range a.Symbols {
		if cl, ok := s.(*CharLiteral); ok {
			cls = append(cls, cl)
		}
	}
	return
}

func (a *Alternate) String() string {
	strs := make([]string, len(a.Symbols))
	for i, s := range a.Symbols {
		strs[i] = s.String()
	}
	return strings.Join(strs, " ")
}

func (a *Alternate) Empty() bool {
	return len(a.Symbols) == 0
}

type Symbol interface {
	IsSymbol()
	Equal(Symbol) bool
	String() string
	GetBSR() bsr.BSR
	GetSymbol() string
}

func (NonTerminal) IsSymbol() {}

type NonTerminal struct {
	bsr bsr.BSR
	NT  string
}

func NewNonTerminal(b bsr.BSR, s string) *NonTerminal {
	return &NonTerminal{
		bsr: b,
		NT:  s,
	}
}

func (nt *NonTerminal) Equal(s Symbol) bool {
	nt1, ok := s.(*NonTerminal)
	if !ok {
		return false
	}
	return nt.NT == nt1.NT
}

func (nt *NonTerminal) GetBSR() bsr.BSR {
	return nt.bsr
}

func (nt *NonTerminal) GetSymbol() string {
	return nt.NT
}

func (nt *NonTerminal) String() string {
	return string(nt.NT)
}

type Terminal interface {
	Symbol
}

func (*Any) IsSymbol()         {}
func (*AnyOf) IsSymbol()       {}
func (*Letter) IsSymbol()      {}
func (*Number) IsSymbol()      {}
func (*Space) IsSymbol()       {}
func (*Upcase) IsSymbol()      {}
func (*Lowcase) IsSymbol()     {}
func (*Not) IsSymbol()         {}
func (*CharLiteral) IsSymbol() {}

func (g *Grammar) IsCharLiteral(symbol string) bool {
	for _, cl := range g.CharLiterals {
		if cl == symbol {
			return true
		}
	}
	return false
}

func (g *Grammar) GetTerminals() []string {
	ts := make([]string, 0, len(g.Terminals))
	for t := range g.Terminals {
		ts = append(ts, t)
	}
	sort.Strings(ts)
	return ts
}

func (g *Grammar) IsTerminal(symbol string) bool {
	// fmt.Printf("ast.IsTerminal(%s)=%t\n", symbol, g.Terminals[symbol])
	return g.Terminals[symbol]
}

func (a *Any) Equal(s Symbol) bool {
	_, ok := s.(*Any)
	return ok
}

func (a *AnyOf) Equal(s Symbol) bool {
	a1, ok := s.(*AnyOf)
	if !ok {
		return false
	}
	return a.Set.Equal(a1.Set)
}

func (l *Letter) Equal(s Symbol) bool {
	_, ok := s.(*Letter)
	return ok
}

func (n *Number) Equal(s Symbol) bool {
	_, ok := s.(*Number)
	return ok
}

func (s *Space) Equal(s1 Symbol) bool {
	_, ok := s1.(*Space)
	return ok
}

func (u *Upcase) Equal(s Symbol) bool {
	_, ok := s.(*Upcase)
	return ok
}

func (l *Lowcase) Equal(s Symbol) bool {
	_, ok := s.(*Lowcase)
	return ok
}

func (n *Not) Equal(s Symbol) bool {
	n1, ok := s.(*Not)
	if !ok {
		return false
	}
	return n.Set.Equal(n1.Set)
}

func (c *CharLiteral) Equal(s Symbol) bool {
	c1, ok := s.(*CharLiteral)
	if !ok {
		return false
	}
	return c.Literal == c1.Literal
}

func (Any) GetSymbol() string {
	return "any"
}

func (a *AnyOf) GetSymbol() string {
	return fmt.Sprintf(`anyof(%s)`, a.Set.JoinEscaped())
}

func (Letter) GetSymbol() string {
	return "letter"
}

func (Number) GetSymbol() string {
	return "number"
}

func (Space) GetSymbol() string {
	return "space"
}

func (Upcase) GetSymbol() string {
	return "upcase"
}

func (Lowcase) GetSymbol() string {
	return "lowcase"
}

func (n *Not) GetSymbol() string {
	return fmt.Sprintf(`not(%s)`, n.Set.JoinEscaped())
}

func (Any) String() string {
	return "any"
}

func (a *AnyOf) String() string {
	return fmt.Sprintf(`anyof "%s"`, strings.Join(a.Set.Elements(), ""))
}

func (Letter) String() string {
	return "letter"
}

func (Number) String() string {
	return "number"
}

func (Space) String() string {
	return "space"
}

func (Upcase) String() string {
	return "upcase"
}

func (Lowcase) String() string {
	return "lowcase"
}

func (*Not) String() string {
	return "not"
}

type Any struct {
	bsr bsr.BSR
}

func NewAny(b bsr.BSR) *Any {
	return &Any{b}
}

func (a *Any) GetBSR() bsr.BSR {
	return a.bsr
}

/*** AnyOf ***/

type AnyOf struct {
	bsr bsr.BSR
	Set *stringset.StringSet
}

func (g *Grammar) NewAnyOf(b bsr.BSR, s string) *AnyOf {
	a := &AnyOf{
		bsr: b,
		Set: stringset.New(newString(s)...),
	}
	g.Terminals[a.GetSymbol()] = true
	return a
}

func (a *AnyOf) GetBSR() bsr.BSR {
	return a.bsr
}

type Letter struct {
	bsr bsr.BSR
}

func NewLetter(b bsr.BSR) *Letter {
	return &Letter{b}
}

func (s *Letter) GetBSR() bsr.BSR {
	return s.bsr
}

type Number struct {
	bsr bsr.BSR
}

func NewNumber(b bsr.BSR) *Number {
	return &Number{b}
}

func (s *Number) GetBSR() bsr.BSR {
	return s.bsr
}

type Space struct {
	bsr bsr.BSR
}

func NewSpace(b bsr.BSR) *Space {
	return &Space{b}
}

func (s *Space) GetBSR() bsr.BSR {
	return s.bsr
}

type Upcase struct {
	bsr bsr.BSR
}

func NewUpcase(b bsr.BSR) *Upcase {
	return &Upcase{b}
}

func (s *Upcase) GetBSR() bsr.BSR {
	return s.bsr
}

type Lowcase struct {
	bsr bsr.BSR
}

func NewLowcase(b bsr.BSR) *Lowcase {
	return &Lowcase{b}
}

func (s *Lowcase) GetBSR() bsr.BSR {
	return s.bsr
}

type Not struct {
	bsr bsr.BSR
	Set *stringset.StringSet
}

func (g *Grammar) NewNot(b bsr.BSR, s string) *Not {
	n := &Not{
		bsr: b,
		Set: stringset.New(newString(s)...),
	}
	g.Terminals[n.GetSymbol()] = true
	return n
}

func (n *Not) GetBSR() bsr.BSR {
	return n.bsr
}

type CharLiteral struct {
	bsr     bsr.BSR
	Literal string
}

func (g *Grammar) NewCharLiteral(b bsr.BSR, s string) *CharLiteral {
	return g.newCharLiteral(b, s[1:len(s)-1])
}

func (c *CharLiteral) GetBSR() bsr.BSR {
	return c.bsr
}

func (c *CharLiteral) GetSymbol() string {
	return fmt.Sprintf(`%s`, c.Literal)
}

func (c *CharLiteral) String() string {
	return fmt.Sprintf(`'%s'`, c.Literal)
}

func newString(s string) (cs []string) {
	rdr := strings.NewReader(s)
	str := ""
	for rdr.Len() > 0 {
		r, _, err := rdr.ReadRune()
		if err != nil {
			panic(err)
		}
		if r == '\\' {
			r, _, err = rdr.ReadRune()
			if err != nil {
				panic(err)
			}
			str = fmt.Sprintf("\\%c", r)
		} else {
			str = string(r)
		}
		cs = append(cs, str)
	}
	return
}

func (g *Grammar) NewString(b bsr.BSR, s string) (cs []Terminal) {
	for _, str := range newString(s[1 : len(s)-1]) {
		cs = append(cs, g.newCharLiteral(b, str))
	}
	return
}

func (g *Grammar) newCharLiteral(b bsr.BSR, s string) *CharLiteral {
	g.CharLiterals = append(g.CharLiterals, s)
	g.Terminals[s] = true
	return &CharLiteral{
		bsr:     b,
		Literal: s,
	}
}
