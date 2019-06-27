package ast

import (
	"fmt"
	"gogll/goutil/utf8"
	"gogll/token"

	"gogll/goutil/stringset"
	"gogll/goutil/stringslice"
)

const (
	Empty = "empty"
)

func GetAST(parseTree interface{}) *Grammar {
	g := parseTree.(*Grammar)
	startSymbol = g.Rules[0].Head.value
	return g
}

type Grammar struct {
	Package *Package

	Rules Rules

	// Key=symbol, Value is first set of symbol
	firstSets map[string]*stringset.StringSet

	// Key=NonTerminal, Value is follow set of NonTerminal
	followSets map[string]*stringset.StringSet
}

func NewGrammar(pkg, rules interface{}) (*Grammar, error) {
	g := &Grammar{
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

func (g *Grammar) FirstOfString(str []string) *stringset.StringSet {
	// fmt.Printf("g.FirstOfString: %s\n", strings.Join(str, " "))
	if len(str) == 0 {
		return stringset.New(Empty)
	}

	first := stringset.New()
	for _, s := range str {
		fs := g.FirstOfSymbol(s)
		first.AddSet(fs)
		if !fs.Contain(Empty) {
			first.Remove(Empty)
			break
		}
	}
	return first
}

func (g *Grammar) FirstOfSymbol(s string) *stringset.StringSet {
	// println("FirstOfSymbol")
	if g.firstSets == nil {
		g.genFirstSets()
	}

	if f, exist := g.firstSets[s]; exist {
		return f
	} else {
		return stringset.New()
	}
}

func (g *Grammar) Follow(nt string) *stringset.StringSet {
	if g.followSets == nil {
		g.genFollow()
	}
	if f, exist := g.followSets[nt]; exist {
		return f
	} else {
		return stringset.New()
	}
}

/*
Dragon book FIRST set algorithm used
*/
func (g *Grammar) genFirstSets() {
	// println("genFirstSets")
	g.initFirstSets()
	for again := true; again; {
		// println(" again")
		again = false
		for _, s := range GetSymbols() {
			// println(" ", s)
			fs := g.getFirstOfSymbol(s)
			if !g.firstSets[s].Equal(fs) {
				g.firstSets[s] = fs
				again = true
			}
		}
		// dumpFirstSets(g.firstSets)
	}
}

func (g *Grammar) initFirstSets() {
	g.firstSets = make(map[string]*stringset.StringSet)
	for _, s := range GetSymbols() {
		g.firstSets[s] = stringset.New()
	}
}

func (g *Grammar) getFirstOfSymbol(s string) *stringset.StringSet {
	// fmt.Println("getFirstOfSymbol: ", s)
	if IsTerminal(s) {
		// fmt.Println("  T: ", stringset.New(s))
		return stringset.New(s)
	}
	// fmt.Println("  NT", g.getFirstOfNonTerminal(s))
	return g.getFirstOfNonTerminal(s)
}

func (g *Grammar) getFirstOfAlternate(a *Alternate) *stringset.StringSet {
	if a.Empty() {
		return stringset.New(Empty)
	}
	return g.FirstOfString(a.Symbols())
}

func (g *Grammar) getFirstOfNonTerminal(s string) *stringset.StringSet {
	first := stringset.New()
	for _, a := range GetRule(s).Alternates {
		f := g.getFirstOfAlternate(a)
		first.Add(f.Elements()...)
	}
	return first
}

/*
Dragon book algoritm used for Follow
*/
func (g *Grammar) genFollow() {
	g.initFollowSets()
	for again := true; again; {
		again = false
		for _, nt := range GetNonTerminals() {
			f := g.genFollowOf(nt)
			if !f.Equal(g.followSets[nt]) {
				again = true
				g.followSets[nt] = f
			}
		}
	}
}

/*
TODO: genFollow only processes syntax rules
*/
func (g *Grammar) genFollowOf(nt string) *stringset.StringSet {
	follow := stringset.New()
	for _, r := range g.Rules {
		for _, a := range r.Alternates {
			bs := stringslice.StringSlice(a.Symbols())
			for _, idx := range bs.Find(nt) {
				first := g.FirstOfString(bs[idx+1:])
				follow.AddSet(first)
				if first.Contain(Empty) {
					follow.AddSet(g.Follow(r.Head.Value()))
				}
			}
		}
	}
	follow.Remove(Empty)
	return follow
}

func (g *Grammar) initFollowSets() {
	g.followSets = make(map[string]*stringset.StringSet)
	for _, nt := range GetNonTerminals() {
		if nt == GetStartSymbol() {
			g.followSets[nt] = stringset.New("$")
		} else {
			g.followSets[nt] = stringset.New()
		}
	}
}

func (g *Grammar) GetRule(head string) *Rule {
	for _, r := range g.Rules {
		if r.Head.Value() == head {
			return r
		}
	}
	return nil
}

// type Action struct {
// 	Tok   *token.Token
// 	Value string
// }

// func NewAction(act interface{}) (*Action, error) {
// 	tok := act.(*token.Token)
// 	a := &Action{
// 		Tok:   tok,
// 		Value: string(tok.Lit[2 : len(tok.Lit)-2]),
// 	}
// 	return a, nil
// }

type Package struct {
	Token *token.Token
}

func NewPackage(pkg interface{}) (*Package, error) {
	tok := pkg.(*token.Token)
	p := &Package{
		Token: tok,
	}
	parserPackage = p.Value()
	return p, nil
}

func (p *Package) Value() string {
	return string(p.Token.Lit)
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

type Head struct {
	Token *token.Token
	value string
}

func NewHead(head interface{}) (*Head, error) {
	tok := head.(*token.Token)
	h := &Head{
		Token: tok,
		value: tok.IDValue(),
	}
	return h, nil
}

func (h *Head) Value() string {
	return h.value
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
	Value() string
	Equal(Symbol) bool
	IsTerminal() bool
	Symbols() Symbols
}

func (*Head) isSymbol()        {}
func (*ID) isSymbol()          {}
func (*AnyChar) isSymbol()     {}
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
		return sym.Value() == sym1.Value()
	}
}
func (sym *ID) Equal(sym1 Symbol) bool {
	if _, ok := sym1.(*ID); !ok {
		return false
	} else {
		return sym.Value() == sym1.Value()
	}
}
func (sym *AnyChar) Equal(sym1 Symbol) bool {
	if _, ok := sym1.(*AnyChar); !ok {
		return false
	} else {
		return sym.Value() == sym1.Value()
	}
}
func (sym *NotString) Equal(sym1 Symbol) bool {
	if _, ok := sym1.(*NotString); !ok {
		return false
	} else {
		return sym.Value() == sym1.Value()
	}
}
func (sym *Space) Equal(sym1 Symbol) bool {
	if _, ok := sym1.(*Space); !ok {
		return false
	} else {
		return sym.Value() == sym1.Value()
	}
}
func (sym *CharLiteral) Equal(sym1 Symbol) bool {
	if _, ok := sym1.(*CharLiteral); !ok {
		return false
	} else {
		return sym.Value() == sym1.Value()
	}
}
func (sym *UpCase) Equal(sym1 Symbol) bool {
	if _, ok := sym1.(*UpCase); !ok {
		return false
	} else {
		return sym.Value() == sym1.Value()
	}
}
func (sym *LowCase) Equal(sym1 Symbol) bool {
	if _, ok := sym1.(*LowCase); !ok {
		return false
	} else {
		return sym.Value() == sym1.Value()
	}
}
func (sym *Letter) Equal(sym1 Symbol) bool {
	if _, ok := sym1.(*Letter); !ok {
		return false
	} else {
		return sym.Value() == sym1.Value()
	}
}
func (sym *Number) Equal(sym1 Symbol) bool {
	if _, ok := sym1.(*Number); !ok {
		return false
	} else {
		return sym.Value() == sym1.Value()
	}
}
func (sym *String) Equal(sym1 Symbol) bool {
	if _, ok := sym1.(*String); !ok {
		return false
	} else {
		return sym.Value() == sym1.Value()
	}
}

func (sym *StringChar) Equal(sym1 Symbol) bool {
	if _, ok := sym1.(*StringChar); !ok {
		return false
	} else {
		return sym.Value() == sym1.Value()
	}
}

func (sym *Head) Symbols() Symbols        { return Symbols{sym.Value()} }
func (sym *ID) Symbols() Symbols          { return Symbols{sym.Value()} }
func (sym *AnyChar) Symbols() Symbols     { return Symbols{sym.Value()} }
func (sym *NotString) Symbols() Symbols   { return Symbols{sym.Value()} }
func (sym *Space) Symbols() Symbols       { return Symbols{sym.Value()} }
func (sym *CharLiteral) Symbols() Symbols { return Symbols{sym.Value()} }
func (sym *UpCase) Symbols() Symbols      { return Symbols{sym.Value()} }
func (sym *LowCase) Symbols() Symbols     { return Symbols{sym.Value()} }
func (sym *Letter) Symbols() Symbols      { return Symbols{sym.Value()} }
func (sym *Number) Symbols() Symbols      { return Symbols{sym.Value()} }
func (s *String) Symbols() Symbols {
	return Symbols(utf8.DecodeRunes(s.Token.Lit[1 : len(s.Token.Lit)-1]))
}
func (sym *StringChar) Symbols() Symbols { return Symbols{sym.Value()} }

type NotString struct {
	Token *token.Token
	value string
}

func (ns *NotString) Value() string {
	return ns.value
}

func NewNotString(sym interface{}) (*NotString, error) {
	tok := sym.(*token.Token)
	symbol := &NotString{
		Token: tok,
		value: "not " + tok.IDValue(),
	}
	return symbol, nil
}

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

func (id *ID) Value() string {
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

func (t *AnyChar) Value() string {
	return t.value
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

func (t *Space) Value() string {
	return t.value
}

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

func (t *CharLiteral) Value() string {
	return t.value
}

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

func (t *UpCase) Value() string {
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

func (t *LowCase) Value() string {
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

func (t *Letter) Value() string {
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

func (t *Number) Value() string {
	return t.value
}

type String struct {
	Token *token.Token
}

func NewString(str interface{}) (*String, error) {
	tok := str.(*token.Token)
	sym := &String{
		Token: tok,
	}
	return sym, nil
}

func (s *String) Value() string {
	return string(s.Token.Lit)
}

type StringChar struct {
	String *String
	value  string
}

func newStringChar(sc string, str *String) *StringChar {
	strCh := &StringChar{
		String: str,
		value:  sc,
	}
	return strCh
}

func (sc *StringChar) Value() string {
	return sc.value
}

/*** Dump ***/

func dumpFirstSets(fs map[string]*stringset.StringSet) {
	for _, s := range GetSymbols() {
		fmt.Printf("%s: %s\n", s, fs[s])
	}
}
