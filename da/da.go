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
	"gogll/goutil/bsr"
	"gogll/parser/symbols"
	"strings"
)

func Go() {
	// fmt.Println("da.Go")
	daReservedWords()
	removeNTsWithoutChildren()
}

func daReservedWords() {
	for _, b := range bsr.GetAll() {
		if reservedWord(b.GetString()) && b.Label.Head() == "NonTerminal" {
			// fmt.Println("daReservedWords", b, b.GetString())
			b.Ignore()
		}
	}
}

// func removeNTsWithoutChildren() {
// 	reps := 0
// 	for again := true; again; {
// 		again = false
// 		for _, b := range bsr.GetAll() {
// 			for i, s := range b.Label.Symbols() {
// 				reps++
// 				if symbols.IsNonTerminal(s) && len(b.GetNTChildrenI(i)) == 0 {
// 					// fmt.Printf("remove %s\n", b)
// 					b.Ignore()
// 					again = true
// 				}
// 			}
// 		}
// 	}
// 	fmt.Printf("da.removeNTsWithoutChildren: %d reps\n", reps)
// }

func removeNTsWithoutChildren() {
	// fmt.Println("da.removeNTsWithoutChildren")
	for _, rt := range bsr.GetRoots() {
		removeZombieChildren(rt)
	}
}

func removeZombieChildren(nt bsr.BSR) {
	// fmt.Println("da.removeZombieChildren", nt)
	if nt.Label.Head() == "Sep" || nt.Label.Head() == "SepE" {
		return
	}

	for i, s := range nt.Label.Symbols() {
		if symbols.IsNonTerminal(s) {
			for _, c := range nt.GetNTChildrenI(i) {
				removeZombieChildren(c)
			}
			if len(nt.GetNTChildrenI(i)) == 0 {
				nt.Ignore()
			}
		}
	}
}

func reservedWord(s string) bool {
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
