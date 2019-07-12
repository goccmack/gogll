package frstflw

import (
	"gogll/ast"
	"gogll/goutil/stringset"
	"gogll/goutil/stringslice"
)

var (
	// Key=symbol, Value is first set of symbol
	firstSets map[string]*stringset.StringSet

	// Key=NonTerminal, Value is follow set of NonTerminal
	followSets map[string]*stringset.StringSet
)

func FirstOfString(str []string) *stringset.StringSet {
	// fmt.Printf("FirstOfString: %s\n", strings.Join(str, " "))
	if len(str) == 0 {
		return stringset.New(ast.Empty)
	}

	first := stringset.New()
	for _, s := range str {
		fs := FirstOfSymbol(s)
		first.AddSet(fs)
		if !fs.Contain(ast.Empty) {
			first.Remove(ast.Empty)
			break
		}
	}
	// fmt.Printf("FirstOfString(%s): %s\n", str, first)
	return first
}

func FirstOfSymbol(s string) *stringset.StringSet {
	// println("FirstOfSymbol")
	if firstSets == nil {
		genFirstSets()
	}

	if f, exist := firstSets[s]; exist {
		return f
	} else {
		return stringset.New()
	}
}

func Follow(nt string) *stringset.StringSet {
	if followSets == nil {
		genFollow()
	}
	if f, exist := followSets[nt]; exist {
		return f
	} else {
		return stringset.New()
	}
}

/*
Dragon book FIRST set algorithm used
*/
func genFirstSets() {
	// println("genFirstSets")
	initFirstSets()
	for again := true; again; {
		// println(" again")
		again = false
		for _, s := range ast.GetSymbols() {
			// println(" ", s)
			fs := getFirstOfSymbol(s)
			if !firstSets[s].Equal(fs) {
				firstSets[s] = fs
				again = true
			}
		}
		// dumpFirstSets(firstSets)
	}
}

func initFirstSets() {
	firstSets = make(map[string]*stringset.StringSet)
	for _, s := range ast.GetSymbols() {
		firstSets[s] = stringset.New()
	}
}

func getFirstOfSymbol(s string) *stringset.StringSet {
	// fmt.Println("getFirstOfSymbol: ", s)
	if ast.IsTerminal(s) {
		// fmt.Println("  T: ", stringset.New(s))
		return stringset.New(s)
	}
	// fmt.Println("  NT", getFirstOfNonTerminal(s))
	return getFirstOfNonTerminal(s)
}

func getFirstOfAlternate(a *ast.Alternate) *stringset.StringSet {
	if a.Empty() {
		return stringset.New(ast.Empty)
	}
	return FirstOfString(a.Symbols())
}

func getFirstOfNonTerminal(s string) *stringset.StringSet {
	first := stringset.New()
	for _, a := range ast.GetRule(s).Alternates {
		f := getFirstOfAlternate(a)
		first.Add(f.Elements()...)
	}
	return first
}

/*
Dragon book algoritm used for Follow
*/
func genFollow() {
	initFollowSets()
	for again := true; again; {
		again = false
		for _, nt := range ast.GetNonTerminals() {
			f := genFollowOf(nt)
			if f.Len() > followSets[nt].Len() {
				again = true
				followSets[nt] = f
			}
		}
	}
}

/*
TODO: genFollow only processes syntax rules
*/
func genFollowOf(nt string) *stringset.StringSet {
	follow := stringset.New()
	for _, r := range ast.GetGrammar().Rules {
		for _, a := range r.Alternates {
			bs := stringslice.StringSlice(a.Symbols())
			for _, idx := range bs.Find(nt) {
				first := FirstOfString(bs[idx+1:])
				follow.AddSet(first)
				if first.Contain(ast.Empty) {
					follow.AddSet(Follow(r.Head.Value()))
				}
			}
		}
	}
	follow.Remove(ast.Empty)
	return follow
}

func initFollowSets() {
	followSets = make(map[string]*stringset.StringSet)
	for _, nt := range ast.GetNonTerminals() {
		if nt == ast.GetStartSymbol() {
			followSets[nt] = stringset.New("$")
		} else {
			followSets[nt] = stringset.New()
		}
	}
}
