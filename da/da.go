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
	"fmt"
	"strings"

	"github.com/goccmack/gogll/goutil/bsr"
	"github.com/goccmack/gogll/parser/symbols"
)

func Go() {
	for _, r := range bsr.GetRoots() {
		da(r, daReservedWord)
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
		if s != "Sep" && symbols.IsNonTerminal(s) {
			if len(b.GetNTChildrenI(i)) > 1 {
				fmt.Printf("  Ambigous: in %s: NT %s (%d) \n", b, s, i)
			}
			for _, b1 := range b.GetNTChildrenI(i) {
				report(b1)
			}
		}
	}
}

func daReservedWord(b bsr.BSR) bool {
	if b.Label.Head() == "NTID" && reservedWord(b.GetString()) {
		// fmt.Printf("daReservedWord: %s\n", b.Label.Head())
		b.Ignore()
		return true
	}
	return false
}

func da(b bsr.BSR, daFuncs ...func(bsr.BSR) bool) {
	if b.Label.Head() == "Sep" || b.Label.Head() == "SepE" || b.Label.Head() == "StringChars" {
		return
	}
	for _, daf := range daFuncs {
		if daf(b) {
			return
		}
	}

	ignore := false
	for i, s := range b.Label.Symbols() {
		if symbols.IsNonTerminal(s) {
			for _, b1 := range b.GetNTChildrenI(i) {
				da(b1, daFuncs...)
			}
			if len(b.GetNTChildrenI(i)) < 1 {
				ignore = true
			}
		}
	}
	if ignore {

		b.Ignore()
	}
}

func reservedWord(s string) bool {
	// fmt.Printf("reservedWord(%s)\n", s)
	switch s {
	case "empty":
		return true
	case "any":
		return true
	// case "anyof":
	// 	return true
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
		// case "not":
		// 	return true
	}
	if strings.HasPrefix(s, "anyof") || strings.HasPrefix(s, "not") {
		return true
	}
	return false
}
