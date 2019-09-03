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
package ast

import (
	"strings"

	"bytes"
	"fmt"
	"gogll/token"

	"gogll/goutil/stringset"
)

const (
	Empty = "empty"
)

func GetAST(parseTree interface{}) *Grammar {
	g := parseTree.(*Grammar)
	// startSymbol = g.Rules[0].Head.stringValue
	return g
}

var (
	g *Grammar
)

type Grammar struct {
	Package *Package

	Rules Rules
}

func GetGrammar() *Grammar {
	return g
}

func NewGrammar(pkg, rules interface{}) (*Grammar, error) {
	g = &Grammar{
		Rules: rules.(Rules),
	}
	if pkg != nil {
		g.Package = pkg.(*Package)
	}
	return g, nil
}

func NewRules(rule interface{}) (Rules, error) {
	rules := Rules{rule.(*Rule)}
	return rules, nil
}

func AddRule(rules, rule interface{}) (Rules, error) {
	rs, r := rules.(Rules), rule.(*Rule)
	rs = append(rs, r)
	return rs, nil
}

type Package struct {
	Token *token.Token
}

func NewPackage(pkg interface{}) (*Package, error) {
	tok := pkg.(*token.Token)
	p := &Package{
		Token: tok,
	}
	if parserPackage != "" {
		err := fmt.Errorf("Duplicate package statement")
		return nil, err
	}
	parserPackage = p.StringValue()
	return p, nil
}

func (p *Package) StringValue() string {
	return p.Token.StringValue()
}

type Alternate struct {
	Body *Body
}

func NewAlternate(body interface{}) (*Alternate, error) {
	a := &Alternate{}
	if body != nil {
		a.Body = body.(*Body)
	}
	return a, nil
}

func (a *Alternate) Empty() bool {
	return a.Body == nil
}

func (a *Alternate) Symbols() Symbols {
	if a.Empty() {
		return Symbols{Empty}
	}
	ss := make(Symbols, 0, len(a.Body.Symbols))
	for _, s := range a.Body.Symbols {
		ss = append(ss, s.Symbols()...)
	}
	return ss
}

type Alternates []*Alternate

func NewAlternates(alt interface{}) (Alternates, error) {
	a := Alternates{alt.(*Alternate)}
	return a, nil
}

func AddAlternate(alts, alt interface{}) (Alternates, error) {
	as := alts.(Alternates)
	as = append(as, alt.(*Alternate))
	return as, nil
}

type Body struct {
	Symbols []Symbol
}

func NewBody(sym interface{}) (*Body, error) {
	sym1 := sym.(Symbol)
	AddSymbol(sym1)
	b := &Body{
		Symbols: []Symbol{sym1},
	}

	return b, nil
}

func AppendSymbol(body, sym interface{}) (*Body, error) {
	b := body.(*Body)
	sym1 := sym.(Symbol)
	AddSymbol(sym1)
	b.Symbols = append(b.Symbols, sym1)
	return b, nil
}

func (b *Body) String() string {
	if b == nil {
		return "ϵ"
	}
	buf := new(bytes.Buffer)
	for i, s := range b.Symbols {
		if i > 0 {
			buf.WriteString(" ")
		}
		buf.WriteString(s.StringValue())
	}
	return buf.String()
}

type Head struct {
	Token       *token.Token
	stringValue string
}

func NewHead(head interface{}) (*Head, error) {
	tok := head.(*token.Token)
	h := &Head{
		Token:       tok,
		stringValue: tok.IDValue(),
	}
	return h, nil
}

func (h *Head) GetPos() token.Pos {
	return h.Token.Pos
}

func (h *Head) StringValue() string {
	return h.stringValue
}

type Rule struct {
	IsStartSymbol bool
	Head          *Head
	Alternates    Alternates
}

func NewRule(start bool, head, alts interface{}) (*Rule, error) {
	r := &Rule{
		IsStartSymbol: start,
		Head:          head.(*Head),
		Alternates:    alts.(Alternates),
	}
	addRule(r)
	return r, nil
}

type Rules []*Rule

type Symbol interface {
	isSymbol()
	GetPos() token.Pos
	StringValue() string
	Equal(Symbol) bool
	IsTerminal() bool
	Symbols() Symbols
}

func (*Head) isSymbol()        {}
func (*ID) isSymbol()          {}
func (*AnyChar) isSymbol()     {}
func (*AnyOf) isSymbol()       {}
func (*NotString) isSymbol()   {}
func (*Space) isSymbol()       {}
func (*String) isSymbol()      {}
func (*CharLiteral) isSymbol() {}
func (*UpCase) isSymbol()      {}
func (*LowCase) isSymbol()     {}
func (*Letter) isSymbol()      {}
func (*Number) isSymbol()      {}
func (*StringChar) isSymbol()  {}

func (*Head) IsTerminal() bool        { return false }
func (*ID) IsTerminal() bool          { return false }
func (*AnyChar) IsTerminal() bool     { return true }
func (*AnyOf) IsTerminal() bool       { return true }
func (*NotString) IsTerminal() bool   { return true }
func (*Space) IsTerminal() bool       { return true }
func (*CharLiteral) IsTerminal() bool { return true }
func (*UpCase) IsTerminal() bool      { return true }
func (*LowCase) IsTerminal() bool     { return true }
func (*Letter) IsTerminal() bool      { return true }
func (*Number) IsTerminal() bool      { return true }
func (*String) IsTerminal() bool      { return true }
func (*StringChar) IsTerminal() bool  { return true }

func (sym *Head) Equal(sym1 Symbol) bool {
	if _, ok := sym1.(*Head); !ok {
		return false
	} else {
		return sym.StringValue() == sym1.StringValue()
	}
}
func (sym *ID) Equal(sym1 Symbol) bool {
	if _, ok := sym1.(*ID); !ok {
		return false
	} else {
		return sym.StringValue() == sym1.StringValue()
	}
}
func (sym *AnyChar) Equal(sym1 Symbol) bool {
	if _, ok := sym1.(*AnyChar); !ok {
		return false
	} else {
		return sym.StringValue() == sym1.StringValue()
	}
}
func (sym *AnyOf) Equal(sym1 Symbol) bool {
	if _, ok := sym1.(*AnyOf); !ok {
		return false
	} else {
		return sym.StringValue() == sym1.StringValue()
	}
}
func (sym *NotString) Equal(sym1 Symbol) bool {
	if _, ok := sym1.(*NotString); !ok {
		return false
	} else {
		return sym.StringValue() == sym1.StringValue()
	}
}
func (sym *Space) Equal(sym1 Symbol) bool {
	if _, ok := sym1.(*Space); !ok {
		return false
	} else {
		return sym.StringValue() == sym1.StringValue()
	}
}
func (sym *CharLiteral) Equal(sym1 Symbol) bool {
	if _, ok := sym1.(*CharLiteral); !ok {
		return false
	} else {
		return sym.StringValue() == sym1.StringValue()
	}
}
func (sym *UpCase) Equal(sym1 Symbol) bool {
	if _, ok := sym1.(*UpCase); !ok {
		return false
	} else {
		return sym.StringValue() == sym1.StringValue()
	}
}
func (sym *LowCase) Equal(sym1 Symbol) bool {
	if _, ok := sym1.(*LowCase); !ok {
		return false
	} else {
		return sym.StringValue() == sym1.StringValue()
	}
}
func (sym *Letter) Equal(sym1 Symbol) bool {
	if _, ok := sym1.(*Letter); !ok {
		return false
	} else {
		return sym.StringValue() == sym1.StringValue()
	}
}
func (sym *Number) Equal(sym1 Symbol) bool {
	if _, ok := sym1.(*Number); !ok {
		return false
	} else {
		return sym.StringValue() == sym1.StringValue()
	}
}
func (sym *String) Equal(sym1 Symbol) bool {
	if _, ok := sym1.(*String); !ok {
		return false
	} else {
		return sym.StringValue() == sym1.StringValue()
	}
}

func (sym *StringChar) Equal(sym1 Symbol) bool {
	if _, ok := sym1.(*StringChar); !ok {
		return false
	} else {
		return sym.StringValue() == sym1.StringValue()
	}
}

func (sym *Head) Symbols() Symbols        { return Symbols{sym.StringValue()} }
func (sym *ID) Symbols() Symbols          { return Symbols{sym.StringValue()} }
func (sym *AnyChar) Symbols() Symbols     { return Symbols{sym.StringValue()} }
func (sym *AnyOf) Symbols() Symbols       { return Symbols{sym.StringValue()} }
func (sym *NotString) Symbols() Symbols   { return Symbols{sym.StringValue()} }
func (sym *Space) Symbols() Symbols       { return Symbols{sym.StringValue()} }
func (sym *CharLiteral) Symbols() Symbols { return Symbols{sym.StringValue()} }
func (sym *UpCase) Symbols() Symbols      { return Symbols{sym.StringValue()} }
func (sym *LowCase) Symbols() Symbols     { return Symbols{sym.StringValue()} }
func (sym *Letter) Symbols() Symbols      { return Symbols{sym.StringValue()} }
func (sym *Number) Symbols() Symbols      { return Symbols{sym.StringValue()} }
func (s *String) Symbols() (ss Symbols)   { return s.symbols }

func (sym *StringChar) Symbols() Symbols { return Symbols{sym.StringValue()} }

type NotString struct {
	Token *token.Token
	value string
}

func (ns *NotString) StringValue() string {
	return ns.value
}

func NewNotString(sym interface{}) (*NotString, error) {
	tok := sym.(*token.Token)
	symbol := &NotString{
		Token: tok,
		// value: "not " + tok.IDValue(),
		value: fmt.Sprintf(`not("%s")`, tok.StringValue()),
	}
	return symbol, nil
}

func (n *NotString) GetPos() token.Pos {
	return n.Token.Pos
}

/*** ID ***/

type ID struct {
	Token *token.Token
	value string
}

func NewID(id interface{}) (*ID, error) {
	tok := id.(*token.Token)
	id1 := &ID{
		Token: tok,
		value: tok.IDValue(),
	}
	return id1, nil
}

func (id *ID) GetPos() token.Pos {
	return id.Token.Pos
}

func (id *ID) StringValue() string {
	return id.value
}

type AnyChar struct {
	Token *token.Token
	value string
}

func NewAnyChar(t interface{}) (*AnyChar, error) {
	tok := t.(*token.Token)
	terminal := &AnyChar{
		Token: tok,
		value: "any",
	}
	return terminal, nil
}

func (t *AnyChar) GetPos() token.Pos {
	return t.Token.Pos
}

func (t *AnyChar) StringValue() string {
	return t.value
}

/*** AnyOf ***/

type AnyOf struct {
	Token *token.Token
	value string
}

func NewAnyOf(str interface{}) (*AnyOf, error) {
	tok := str.(*token.Token)
	anyOf := &AnyOf{
		Token: tok,
		value: fmt.Sprintf(`anyof("%s")`, tok.StringValue()),
	}
	return anyOf, nil
}

func (a *AnyOf) GetPos() token.Pos {
	return a.Token.Pos
}

func (a *AnyOf) StringValue() string {
	return a.value
}

type Space struct {
	Token *token.Token
	value string
}

func NewSpace(t interface{}) (*Space, error) {
	tok := t.(*token.Token)
	terminal := &Space{
		Token: tok,
		value: "space",
	}
	return terminal, nil
}

func (s *Space) GetPos() token.Pos {
	return s.Token.Pos
}

func (t *Space) StringValue() string {
	return t.value
}

/*** CharLiteral ***/

type CharLiteral struct {
	Token *token.Token
	value string
	Rune  rune
}

func NewCharLiteral(t interface{}) (*CharLiteral, error) {
	tok := t.(*token.Token)
	r, err := tok.UTF8Rune()
	if err != nil {
		return nil, err
	}
	terminal := &CharLiteral{
		Token: tok,
		value: tok.CharLiteralValue(),
		Rune:  r,
	}
	return terminal, nil
}

func (t *CharLiteral) StringValue() string {
	return t.value
}

func (c *CharLiteral) GetPos() token.Pos {
	return c.Token.Pos
}

/*** UpCase ***/

type UpCase struct {
	Token *token.Token
	value string
}

func NewUpCase(t interface{}) (*UpCase, error) {
	tok := t.(*token.Token)
	terminal := &UpCase{
		Token: tok,
		value: tok.IDValue(),
	}
	return terminal, nil
}

func (u *UpCase) GetPos() token.Pos {
	return u.Token.Pos
}

func (t *UpCase) StringValue() string {
	return t.value
}

type LowCase struct {
	Token *token.Token
	value string
}

func NewLowCase(t interface{}) (*LowCase, error) {
	tok := t.(*token.Token)
	terminal := &LowCase{
		Token: tok,
		value: tok.IDValue(),
	}
	return terminal, nil
}

func (l *LowCase) GetPos() token.Pos {
	return l.Token.Pos
}

func (t *LowCase) StringValue() string {
	return t.value
}

type Letter struct {
	Token *token.Token
	value string
}

func NewLetter(t interface{}) (*Letter, error) {
	tok := t.(*token.Token)
	terminal := &Letter{
		Token: tok,
		value: tok.IDValue(),
	}
	return terminal, nil
}

func (l *Letter) GetPos() token.Pos {
	return l.Token.Pos
}

func (t *Letter) StringValue() string {
	return t.value
}

type Number struct {
	Token *token.Token
	value string
}

func NewNumber(t interface{}) (*Number, error) {
	tok := t.(*token.Token)
	terminal := &Number{
		Token: tok,
		value: tok.IDValue(),
	}
	return terminal, nil
}

func (n *Number) GetPos() token.Pos {
	return n.Token.Pos
}

func (t *Number) StringValue() string {
	return t.value
}

// String

type String struct {
	Token    *token.Token
	strChars []*StringChar
	symbols  []string
}

func NewString(str interface{}) (*String, error) {
	tok := str.(*token.Token)
	sym := &String{
		Token: tok,
	}
	var err error
	sym.strChars, err = newStringChars(sym)
	if err != nil {
		return nil, err
	}
	for _, sc := range sym.strChars {
		sym.symbols = append(sym.symbols, sc.stringValue)
	}
	return sym, nil
}

func newStringChars(s *String) (cs []*StringChar, err error) {
	rdr := strings.NewReader(s.Token.StringValue())
	str := ""
	for rdr.Len() > 0 {
		r, _, err := rdr.ReadRune()
		if err != nil {
			return nil, err
		}
		if r == '\\' {
			r, _, err = rdr.ReadRune()
			if err != nil {
				return nil, err
			}
			str = fmt.Sprintf("\\%c", r)
		} else {
			str = string(r)
		}
		cs = append(cs, newStringChar(r, str, s))
	}
	return
}

func (s *String) GetPos() token.Pos {
	return s.Token.Pos
}

func (s *String) StringValue() string {
	return string(s.Token.Lit)
}

type StringChar struct {
	String      *String
	Rune        rune
	stringValue string
}

func newStringChar(r rune, sc string, str *String) *StringChar {
	// fmt.Printf("ast.newStringChar*(%s)\n", sc)
	strCh := &StringChar{
		String:      str,
		Rune:        r,
		stringValue: sc,
	}
	return strCh
}

func (s *StringChar) GetPos() token.Pos {
	return s.GetPos()
}

func (sc *StringChar) StringValue() string {
	return sc.stringValue
}

/*** Dump ***/

func dumpFirstSets(fs map[string]*stringset.StringSet) {
	for _, s := range GetSymbols() {
		fmt.Printf("%s: %s\n", s, fs[s])
	}
}
