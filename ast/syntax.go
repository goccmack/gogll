package ast

// Syntax rules

type SyntaxAlternate struct {
	Symbols []SyntaxSymbol
}

type SyntaxRule struct {
	Head       *NT
	Alternates []*SyntaxAlternate
}

type SyntaxSymbol interface {
	isSyntaxSymbol()
	// Lext returns the left extent of SyntaxSymbol in the input string
	Lext() int
	Token() string
	String() string
}

func (*NT) isSyntaxSymbol() {}

// Terminals
func (*TokID) isSyntaxSymbol()     {}
func (*StringLit) isSyntaxSymbol() {}

func (a *SyntaxAlternate) GetSymbols() []string {
	symbols := make([]string, len(a.Symbols))
	for i, s := range a.Symbols {
		symbols[i] = s.Token()
	}
	return symbols
}

func (a *SyntaxAlternate) Empty() bool {
	return len(a.Symbols) == 0
}

// ID returns the head of rule r
func (r *SyntaxRule) ID() string {
	return r.Head.ID()
}

func (r *SyntaxRule) Lext() int {
	return r.Head.Lext()
}
