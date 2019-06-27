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
		checkRuleReference(nt)
	}
}

func checkRuleReference(nt string) {
	if _, exist := ast.GetRules()[nt]; !exist {
		fail("No rule declaration for " + nt)
	}
}
