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

// The lexer part of the AST

package ast

import (
	"bytes"
	"fmt"
	"unicode"

	"github.com/goccmack/gogll/v3/token"
	"github.com/goccmack/gogll/v3/util/runeset"
)

// This file contains the AST components for lexical rules

// TriState has values: {Undefined, False, True}
type TriState int

const (
	// Undefined is a TriState value
	Undefined TriState = iota
	// False is a TriState value
	False
	// True is a TriState value
	True
)

type Any struct {
	tok *token.Token
}

type AnyOf struct {
	any    *token.Token
	strLit *token.Token
	Set    *runeset.RuneSet
}

type CharLiteral struct {
	tok     *token.Token
	Literal []rune
}

type LexBracket struct {
	leftBracket *token.Token
	Type        BracketType
	Alternates  []*RegExp
}

type BracketType int

const (
	LexGroup BracketType = iota
	LexOptional
	LexZeroOrMore
	LexOneOrMore
)

type LexBase interface {
	isLexBase()
	LexSymbol
	Equal(LexBase) bool
}

func (*Any) isLexBase()          {}
func (*AnyOf) isLexBase()        {}
func (*CharLiteral) isLexBase()  {}
func (*Not) isLexBase()          {}
func (*UnicodeClass) isLexBase() {}
func (*UnicodeSet) isLexBase()   {}

type LexRule struct {
	Suppress bool
	TokID    *TokID
	RegExp   *RegExp
}

type LexSymbol interface {
	isLexSymbol()
	Lext() int
	String() string
}

func (*Any) isLexSymbol()          {}
func (*AnyOf) isLexSymbol()        {}
func (*CharLiteral) isLexSymbol()  {}
func (*LexBracket) isLexSymbol()   {}
func (*Not) isLexSymbol()          {}
func (*UnicodeClass) isLexSymbol() {}
func (*UnicodeSet) isLexSymbol()   {}

type Not struct {
	not    *token.Token
	strLit *token.Token
	Set    *runeset.RuneSet
}

type RegExp struct {
	Symbols []LexSymbol
}

type StringLit struct {
	tok *token.Token
}

type UnicodeClass struct {
	tok  *token.Token
	Type UnicodeClassType
}

type UnicodeClassType int

const (
	Letter UnicodeClassType = iota
	Upcase
	Lowcase
	Number
	Space
)

type UnicodeSet struct {
	lext   int
	Pos    *Position
	Ranges UnicodeRanges
}

type UnicodeRanges []*UnicodeRange
type UnicodeRange struct {
	lext    int
	Pos     *Position
	Type    UnicodeRangeType
	Exclude bool
	Range   string
}

type UnicodeRangeType int

const (
	UnicodeCategory UnicodeRangeType = iota
	UnicodeProperty
)

/*** Methods ***/

func (*Any) Equal(other LexBase) bool {
	if other == nil {
		return false
	}
	_, ok := other.(*Any)
	return ok
}

func (a *Any) Lext() int {
	return a.tok.Lext()
}

func (ao *AnyOf) Equal(other LexBase) bool {
	if other == nil {
		return false
	}
	ao1, ok := other.(*AnyOf)
	if !ok {
		return false
	}
	return ao.Set.Equal(ao1.Set)
}

func (a *AnyOf) Lext() int {
	return a.any.Lext()
}

func NewCharLiteral(tok *token.Token, literal []rune) *CharLiteral {
	return &CharLiteral{
		tok:     tok,
		Literal: literal,
	}
}

func (c *CharLiteral) Char() rune {
	if c.Literal[1] == '\\' {
		switch c.Literal[2] {
		case '\'':
			return '\''
		case '"':
			return '"'
		case '\\':
			return '\\'
		case 't':
			return '\t'
		case 'n':
			return '\n'
		case 'r':
			return '\r'
		default:
			panic(fmt.Sprintf("invalid '%c'", c.Literal[2]))
		}
	} else {
		return c.Literal[1]
	}
}

func (c *CharLiteral) Equal(other LexBase) bool {
	if other == nil {
		return false
	}
	c1, ok := other.(*CharLiteral)
	if !ok {
		return false
	}
	// fmt.Printf("'%c'.Equal('%c') = %t\n", c.Char(), c1.Char(), c.Char() == c1.Char())
	return c.Char() == c1.Char()
}

func (c *CharLiteral) Lext() int {
	return c.tok.Lext()
}

func (l *LexBracket) LeftBracket() string {
	switch l.Type {
	case LexGroup:
		return "("
	case LexOptional:
		return "["
	case LexZeroOrMore:
		return "{"
	case LexOneOrMore:
		return "<"
	}
	panic("invalid")
}

func (l *LexBracket) RightBracket() string {
	switch l.Type {
	case LexGroup:
		return ")"
	case LexOptional:
		return "]"
	case LexZeroOrMore:
		return "}"
	case LexOneOrMore:
		return ">"
	}
	panic("invalid")
}

// Returns the id of the lex rule
func (l *LexRule) ID() string {
	return l.TokID.ID()
}

func (l *LexRule) Lext() int {
	return l.TokID.Lext()
}

func (l *LexRule) String() string {
	return fmt.Sprintf("%s : %s ;", l.ID(), l.RegExp)
}

func (b *LexBracket) Lext() int {
	return b.leftBracket.Lext()
}

func (n *Not) Equal(other LexBase) bool {
	if other == nil {
		return false
	}
	n1, ok := other.(*Not)
	if !ok {
		return false
	}
	return n.Set.Equal(n1.Set)
}

func (n *Not) Lext() int {
	return n.not.Lext()
}

func (re *RegExp) String() string {
	w := new(bytes.Buffer)
	for _, symbol := range re.Symbols {
		fmt.Fprint(w, symbol)
	}
	return w.String()
}

func (u *UnicodeClass) Equal(other LexBase) bool {
	if other == nil {
		return false
	}
	u1, ok := other.(*UnicodeClass)
	if !ok {
		return false
	}
	return u.Type == u1.Type
}

func (u *UnicodeClass) Lext() int {
	return u.Lext()
}

func (*Any) String() string {
	return "."
}

func (a *AnyOf) String() string {
	return fmt.Sprintf("any %s", string(a.strLit.Literal()))
}

func (c *CharLiteral) String() string {
	return string(c.Literal)
}

func (lb *LexBracket) String() string {
	w := new(bytes.Buffer)
	fmt.Fprint(w, lb.LeftBracket())
	for i, alt := range lb.Alternates {
		if i > 0 {
			fmt.Fprint(w, " | ")
		}
		fmt.Fprint(w, alt)
	}
	fmt.Fprint(w, lb.RightBracket())
	return w.String()
}

func (n *Not) String() string {
	return fmt.Sprintf("not %s", string(n.strLit.Literal()))
}

func (sl *StringLit) ContainsWhiteSpace() bool {
	for _, r := range sl.tok.LiteralStripEscape() {
		switch r {
		case ' ', '\t', '\n', '\r':
			return true
		}
	}
	return false
}

func (sl *StringLit) ID() string {
	return string(sl.Value())
}

func (sl *StringLit) Literal() []rune {
	return sl.tok.Literal()
}

func (sl *StringLit) Value() []rune {
	slit := sl.tok.LiteralStripEscape()
	value := slit[1 : len(slit)-1]
	// fmt.Printf("*StringLit.Value %s %s\n", string(slit), string(value))
	return value
}

func (u *UnicodeClass) ContainsSet(s *UnicodeSet) bool {
	for _, rng := range s.Ranges {
		if !rangeTableInRangeTable(rng.GetRangeTable(), u.GetRangeTable()) {
			return false
		}
	}
	return true
}

func (u *UnicodeClass) GetRangeTable() *unicode.RangeTable {
	switch u.Type {
	case Letter:
		return unicode.Letter
	case Upcase:
		return unicode.Lu
	case Lowcase:
		return unicode.Ll
	case Number:
		return unicode.Number
	case Space:
		return unicode.Space
	}
	panic("impossible")
}

func (u *UnicodeClass) String() string {
	return string(u.tok.Literal())
}

func (u *UnicodeRange) Equals(u1 *UnicodeRange) bool {
	return u.Type == u1.Type &&
		u.Exclude == u1.Exclude &&
		u.Range == u1.Range
}

// EqualsRange ignores the Exclude field
func (u *UnicodeRange) EqualsRange(u1 *UnicodeRange) bool {
	return u.Type == u1.Type &&
		u.Range == u1.Range
}

func (u *UnicodeRange) String() string {
	return u.Range
}

func (u UnicodeRanges) Contain(rng *UnicodeRange) bool {
	for _, rng1 := range u {
		if rng1.Equals(rng) {
			return true
		}
	}
	return false
}

// ContainRanges ignores the Exclude field
func (u UnicodeRanges) ContainRange(rng *UnicodeRange) bool {
	for _, rng1 := range u {
		if rng1.EqualsRange(rng) {
			return true
		}
	}
	return false
}

func (u UnicodeRanges) Equal(other UnicodeRanges) bool {
	if other == nil {
		return false
	}
	for _, rng := range u {
		if !other.Contain(rng) {
			return false
		}
	}
	return true
}

// EqualRange ignores Exclude field
func (u UnicodeRanges) EqualRange(other UnicodeRanges) bool {
	if other == nil {
		return false
	}
	for _, rng := range u {
		if !other.ContainRange(rng) {
			return false
		}
	}
	return true
}

func (u UnicodeRanges) Lext() int {
	return u[0].lext
}

func (u UnicodeRanges) String() string {
	w := new(bytes.Buffer)
	for i, rng := range u {
		if i > 0 {
			fmt.Fprint(w, ",")
		}
		if rng.Exclude {
			fmt.Fprint(w, "-")
		}
		fmt.Fprintf(w, "\\p{%s}", rng)
	}
	return w.String()
}

func (u *UnicodeRange) ContainsClass(c *UnicodeClass) bool {
	switch c.Type {
	case Letter:
		return rangeTableInRangeTable(unicode.Letter, u.GetRangeTable())
	case Upcase:
		return rangeTableInRangeTable(unicode.Lu, u.GetRangeTable())
	case Lowcase:
		return rangeTableInRangeTable(unicode.Ll, u.GetRangeTable())
	case Number:
		return rangeTableInRangeTable(unicode.Number, u.GetRangeTable())
	case Space:
		return rangeTableInRangeTable(unicode.Space, u.GetRangeTable())
	}
	panic("impossible")
}

// rangeInRange returns true iff b contains a
func rangeTableInRangeTable(a, b *unicode.RangeTable) bool {
	if a == nil || b == nil {
		fmt.Printf("a=%p, b=%p\n", a, b)
		panic("nil")
	}
	for _, rng := range a.R32 {
		if !unicode.In(rune(rng.Lo), b) || !unicode.In(rune(rng.Hi), b) {
			return false
		}
	}
	return true
}

// GetRangeTable returns the category or property RangeTable of r
func (u *UnicodeRange) GetRangeTable() *unicode.RangeTable {
	switch u.Type {
	case UnicodeCategory:
		if unicode.Categories[u.Range] == nil {
			panic(u.Range)
		}
		return unicode.Categories[u.Range]
	case UnicodeProperty:
		if unicode.Properties[u.Range] == nil {
			panic(u.Range)
		}
		return unicode.Properties[u.Range]
	}
	panic("impossible")
}

// In returns true iff rng contains u
func (u *UnicodeRange) In(rng *UnicodeRange) bool {
	return rangeTableInRangeTable(u.GetRangeTable(), rng.GetRangeTable())
}

// Subset returns true iff f is a subset or equal of sup
func (u *UnicodeRange) Subset(sup *UnicodeSet) bool {
	for _, rsup := range sup.Ranges {
		if u.In(rsup) {
			return true
		}
	}
	return false
}

func (u *UnicodeSet) ContainsClass(c *UnicodeClass) bool {
	for _, rng := range u.Ranges {
		if rng.ContainsClass(c) {
			return true
		}
	}
	return false
}

func (u *UnicodeSet) ContainsRune(r rune) bool {
	incl, excl := u.GetRangeTables()
	return unicode.In(r, incl...) && !unicode.In(r, excl...)
}

func (u *UnicodeSet) ContainsRunes(runes ...rune) bool {
	for _, r := range runes {
		if !u.ContainsRune(r) {
			return false
		}
	}
	return true
}

func (u *UnicodeSet) ContainsRangeTable(rngTab *unicode.RangeTable) bool {
	for _, rng := range u.Ranges {
		if rangeTableInRangeTable(rngTab, rng.GetRangeTable()) {
			return true
		}
	}
	return false
}

func (u *UnicodeSet) ContainsSet(u1 *UnicodeSet) bool {
	for _, rng := range u1.Ranges {
		if !u.ContainsRangeTable(rng.GetRangeTable()) {
			return false
		}
	}
	return true
}

func (u *UnicodeSet) Equal(other LexBase) bool {
	if other == nil {
		return false
	}
	u1, ok := other.(*UnicodeSet)
	if !ok {
		return false
	}
	return u.Ranges.Equal(u1.Ranges)
}

func (u *UnicodeSet) GetRangeTables() (incl, excl []*unicode.RangeTable) {
	for _, rng := range u.Ranges {
		if rng.Exclude {
			excl = append(excl, rng.GetRangeTable())
		} else {
			incl = append(incl, rng.GetRangeTable())
		}
	}
	return
}

func (u *UnicodeSet) Lext() int {
	return u.lext
}

func (u *UnicodeSet) String() string {
	return "'[" + u.Ranges.String() + "]'"
}

// Subset returns true iff u is a subset or equal of sup
func (u *UnicodeSet) Subset(sup *UnicodeSet) bool {
	for _, rng := range u.Ranges {
		if !rng.Subset(sup) {
			return false
		}
	}
	return true
}

/*** Utils ***/

// StringLitToTokID returns a dummy TokID with ID = id
func StringLitToTokID(id *StringLit) *TokID {
	return &TokID{
		token.New(token.StringToType["tokid"],
			id.tok.Lext()+1, id.tok.Rext()-1, id.tok.GetInput()),
	}
}

// CharLitFromStringLit returns a dummy CharLiteral with Literal sl.Literal[i]
// If escaped sl.Literal[i] == '\\' and sl.Literal[i+1] is the escaped char.
func CharLitFromStringLit(sl *StringLit, i int, escaped bool) *CharLiteral {
	// Make char literal
	lit := []rune{'\''}
	if escaped {
		if sl.Literal()[i+1] != '"' {
			lit = append(lit, '\\')
		}
		lit = append(lit, sl.Literal()[i+1])
	} else {
		lit = append(lit, sl.Literal()[i])
	}
	lit = append(lit, '\'')

	rext := sl.Lext() + i + 1
	if escaped {
		rext++
	}

	cl := NewCharLiteral(
		token.New(
			token.StringToType["char_lit"],
			sl.Lext()+i, rext, sl.tok.GetInput()),
		lit)
	return cl
}
