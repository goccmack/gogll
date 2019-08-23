
package symbols

func IsNonTerminal(symbol string) bool {
	return nonTerminals[symbol]
}

func IsTerminal(symbol string) bool {
	return !nonTerminals[symbol]
}

var nonTerminals = map[string]bool{ 
	"A":true,
	"B":true,
	"S":true,
}
