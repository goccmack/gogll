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