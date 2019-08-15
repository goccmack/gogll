/*
Package check performs semantic checks on the grammar.
*/
package check

import "gogll/ast"

func Check(g *ast.Grammar) {
	checkRuleReferences()
}

func checkRuleReferences() {
	for _, nt := range ast.GetNonTerminals() {
		checkRuleReference(ast.GetSymbol(nt))
	}
}

func checkRuleReference(nt ast.Symbol) {
	if _, exist := ast.GetRules()[nt.StringValue()]; !exist {
		fail(nt.GetPos(), "No rule declaration for %s", nt.StringValue())
	}
}
