
package symbols

func IsNonTerminal(symbol string) bool {
	return nonTerminals[symbol]
}

func IsTerminal(symbol string) bool {
	return !nonTerminals[symbol]
}

var nonTerminals = map[string]bool{ 
	"Alternate":true,
	"Alternates":true,
	"CharLiteral":true,
	"EscapedChar":true,
	"GoGLL":true,
	"Head":true,
	"NTChar":true,
	"NTChars":true,
	"NTID":true,
	"NonTerminal":true,
	"Package":true,
	"Rule":true,
	"Rules":true,
	"Sep":true,
	"SepChar":true,
	"SepE":true,
	"String":true,
	"StringChars":true,
	"Symbol":true,
	"Symbols":true,
	"Terminal":true,
}
