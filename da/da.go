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
Package da disambiguates the BSR set.
*/
package da

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/goccmack/gogll/goutil/bsr"
	"github.com/goccmack/gogll/parser"
	"github.com/goccmack/gogll/parser/symbols"
)

func Go() {
	for _, r := range bsr.GetRoots() {
		da(r)
	}
}

// Report lists the ambiguous subtrees of the parse forest
func Report() {
	rts := bsr.GetRoots()
	if len(rts) != 1 {
		fmt.Println(len(rts), "ambiguous BSR roots")
	}
	for i, b := range bsr.GetRoots() {
		fmt.Println("Root", i)
		report(b)
	}
}

func report(b bsr.BSR) {
	for i, s := range b.Label.Symbols() {
		ln, col := parser.GetLineColumn(b.LeftExtent())
		if s != "Sep" && symbols.IsNonTerminal(s) {
			if len(b.GetNTChildrenI(i)) != 1 {
				fmt.Printf("  Ambigous: in %s: NT %s (%d) at line %d col %d \n", b, s, i, ln, col)
				fmt.Println("   Children:")
				for _, c := range b.GetNTChildrenI(i) {
					fmt.Printf("     %s\n", c)
				}
			}
			for _, b1 := range b.GetNTChildrenI(i) {
				report(b1)
			}
		}
	}
}

func daReservedWord(b bsr.BSR) bool {
	if !reservedWord(b.GetString()) {
		return false
	}
	if b.Label.Head() == "NTID" || b.Label.Head() == "NonTerminal" {
		// ln, col := parser.GetLineColumn(b.LeftExtent())
		// fmt.Printf("Ignore: daReservedWord: %s at line %d col %d\n", b.Label.Head(), ln, col)
		b.Ignore()
		return true
	}
	return false
}

func indent(lvl int) string {
	buf := new(bytes.Buffer)
	for i := 0; i < lvl*2; i++ {
		buf.WriteString(" ")
	}
	return buf.String()
}

func da(b bsr.BSR) {
	// ln, col := parser.GetLineColumn(b.LeftExtent())
	// fmt.Printf("da: %s at line %d col %d\n", b, ln, col)

	for i, s := range b.Label.Symbols() {
		if s != "Sep" && symbols.IsNonTerminal(s) {
			for _, c := range b.GetNTChildrenI(i) {
				da(c)
			}
		}
	}
	switch b.Label.Head() {
	case "Symbol":
		daSymbol(b)
	case "Symbols":
		daSymbols(b)
	case "Alternate":
		daAlternate(b)
		// case "Alternates":
		// 	daAlternates(b)
		// case "Rule":
		// 	daRule(b)
		// case "Rules":
		// 	daRules(b)
	}
}

/*
Alternate
    :   Symbols
    |   "empty"
    ;
*/
func daAlternate(b bsr.BSR) {
	if b.Alternate() == 0 && len(b.GetNTChildrenI(0)) < 1 {
		b.Ignore()
	}
}

/*
Alternates
    :   Alternate
    |   Alternate SepE "|" SepE Alternates
    ;
*/
func daAlternates(b bsr.BSR) {
	if len(b.GetNTChildrenI(0)) < 1 ||
		b.Alternate() == 1 && len(b.GetNTChildrenI(4)) < 1 {
		b.Ignore()
	}
}

// Symbol : NonTerminal | Terminal ;
func daSymbol(b bsr.BSR) {
	if b.Alternate() == 0 && startsWith(b, resWords...) {
		b.Ignore()
	}
}

/*
Rule : Head SepE ":" SepE Alternates SepE ";" ;
*/
func daRule(b bsr.BSR) {
	if len(b.GetNTChildrenI(4)) < 1 {
		b.Ignore()
	}
}

/*
Rules
    :   Rule
    |   Rule SepE Rules
    ;
*/
func daRules(b bsr.BSR) {
	if len(b.GetNTChildrenI(0)) < 1 ||
		b.Alternate() == 1 && len(b.GetNTChildrenI(2)) < 1 {
		b.Ignore()
	}
}

/*
Symbols
    :   Symbol Sep Symbols
    |   Symbol
    ;
*/
func daSymbols(b bsr.BSR) {
	// ln, col := parser.GetLineColumn(b.LeftExtent())
	// fmt.Printf("daSymbols: %s at ln %d, col %d\n", b, ln, col)

	if b.Alternate() == 1 && len(b.GetNTChildrenI(0)) != 0 {
		return
	}
	if len(b.GetNTChildrenI(0)) < 1 || len(b.GetNTChildrenI(2)) < 1 {
		b.Ignore()
		return
	}
	sym := b.GetNTChildI(0).GetString()
	if (sym == "not" || sym == "anyof") && strings.HasPrefix(b.GetNTChildI(2).GetString(), "\"") {
		b.Ignore()
	}
}

// func da(b bsr.BSR, lvl int, daFuncs ...func(bsr.BSR) bool) {
// 	ln, col := parser.GetLineColumn(b.LeftExtent())
// 	fmt.Printf("%sda.da: %s at ln %d col %d\n", indent(lvl), b, ln, col)

// 	if b.Label.Head() == "Sep" || b.Label.Head() == "SepE" || b.Label.Head() == "StringChars" {
// 		return
// 	}

// 	for _, daf := range daFuncs {
// 		if daf(b) {
// 			fmt.Printf("%sIgnore %s\n", indent(lvl), b.Label)
// 			return
// 		}
// 	}

// 	ignore := false
// 	for i, s := range b.Label.Symbols() {
// 		if symbols.IsNonTerminal(s) {
// 			for _, b1 := range b.GetNTChildrenI(i) {
// 				da(b1, lvl+1, daFuncs...)
// 			}
// 			if len(b.GetNTChildrenI(i)) < 1 {
// 				fmt.Printf("%s%sNT %d has 0 children\n", indent(lvl), b.Label, i)
// 				ignore = true
// 			}
// 		}
// 	}
// 	if ignore {
// 		fmt.Printf("%s%s0 children - ignore \n", indent(lvl), b.Label)
// 		b.Ignore()
// 	}
// }

func reservedWord(s string) bool {
	// fmt.Printf("reservedWord(%s)\n", s)
	switch s {
	case "empty":
		return true
	case "any":
		return true
	case "anyof":
		return true
	case "letter":
		return true
	case "number":
		return true
	case "space":
		return true
	case "upcase":
		return true
	case "lowcase":
		return true
	case "not":
		return true
	}
	if strings.HasPrefix(s, "anyof") || strings.HasPrefix(s, "not") {
		return true
	}
	return false
}

var resWords = []string{
	"empty",
	"any",
	"anyof",
	"letter",
	"number",
	"space",
	"upcase",
	"lowcase",
	"not",
}

func startsWith(b bsr.BSR, words ...string) bool {
	for _, w := range words {
		if strings.HasPrefix(b.GetString(), w) {
			return true
		}
	}
	return false
}
