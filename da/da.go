/*
Package da disambiguates the BSR set.
*/
package da

import (
	"gogll/goutil/bsr"
	"gogll/parser/symbols"
)

func Go() {
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
	for _, rt := range bsr.GetRoots() {
		removeZombieChildren(rt)
	}
}

func removeZombieChildren(nt bsr.BSR) {
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
	return false
}
