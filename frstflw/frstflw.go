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

package frstflw

import (
	"github.com/goccmack/gogll/v3/ast"
	"github.com/goccmack/goutil/stringset"
	"github.com/goccmack/goutil/stringslice"
)

const Empty = "ϵ"

type FF struct {
	// Key=symbol, Value is first set of symbol
	firstSets map[string]*stringset.StringSet

	// Key=NonTerminal, Value is follow set of NonTerminal
	followSets map[string]*stringset.StringSet

	g *ast.GoGLL
}

func New(g *ast.GoGLL) *FF {
	ff := &FF{
		g: g,
	}
	ff.genFirstSets()
	ff.genFollow()
	return ff
}

func (ff *FF) FirstOfString(str []string) *stringset.StringSet {
	// fmt.Printf("FirstOfString: %s\n", strings.Join(str, " "))
	if len(str) == 0 {
		return stringset.New(Empty)
	}

	first := stringset.New()
	for _, s := range str {
		fs := ff.FirstOfSymbol(s)
		first.AddSet(fs)
		if !fs.Contain(Empty) {
			first.Remove(Empty)
			break
		}
	}
	// fmt.Printf("FirstOfString(%s): %s\n", strings.Join(str, " "), first)
	return first
}

func (ff *FF) FirstOfSymbol(s string) *stringset.StringSet {
	// fmt.Printf("frstflw.FirstOfSymbol(%s)\n", s)
	if f, exist := ff.firstSets[s]; exist {
		return f
	}
	return stringset.New()
}

func (ff *FF) Follow(nt string) *stringset.StringSet {
	if f, exist := ff.followSets[nt]; exist {
		return f
	} else {
		return stringset.New()
	}
}

/*
Dragon book FIRST set algorithm used
*/
func (ff *FF) genFirstSets() {
	// println("genFirstSets")
	ff.initFirstSets()
	for again := true; again; {
		// println(" again")
		again = false
		for _, s := range ff.g.GetSymbols() {
			// println(" ", s)
			fs := ff.getFirstOfSymbol(s)
			// fmt.Printf("  fs=%s eq=%t\n", fs.Elements(), ff.firstSets[s].Equal(fs))
			if !ff.firstSets[s].Equal(fs) {
				// fmt.Printf(" changed\n")
				ff.firstSets[s] = fs
				again = true
			}
		}
	}
	// for sym, fs := range ff.firstSets {
	// 	fmt.Printf("First(\"%s\"):%s\n", sym, fs.Elements())
	// }
}

func (ff *FF) initFirstSets() {
	ff.firstSets = make(map[string]*stringset.StringSet)
	for _, s := range ff.g.GetSymbols() {
		ff.firstSets[s] = stringset.New()
	}
}

func (ff *FF) getFirstOfSymbol(s string) *stringset.StringSet {
	// fmt.Println("getFirstOfSymbol: ", s)
	if ff.g.Terminals.Contain(s) {
		fst := stringset.New(s)
		// fmt.Println("  T: ", stringset.New(s))
		return fst
	}
	fst := ff.getFirstOfNonTerminal(s)
	// fmt.Println("  NT", fst)
	return fst
}

func (ff *FF) getFirstOfAlternate(a *ast.SyntaxAlternate) *stringset.StringSet {
	if a.Empty() {
		return stringset.New(Empty)
	}
	return ff.FirstOfString(a.GetSymbols())
}

func (ff *FF) getFirstOfNonTerminal(s string) *stringset.StringSet {
	first := stringset.New()
	for _, a := range ff.g.GetSyntaxRule(s).Alternates {
		f := ff.getFirstOfAlternate(a)
		first.Add(f.Elements()...)
	}
	return first
}

/*
Dragon book algoritm used for Follow
*/
func (ff *FF) genFollow() {
	ff.initFollowSets()
	for again := true; again; {
		again = false
		numSets := len(ff.followSets)
		for _, nt := range ff.g.NonTerminals.Elements() {
			f := ff.genFollowOf(nt)
			if f.Len() != ff.followSets[nt].Len() {
				again = true
				ff.followSets[nt] = f
			}
		}
		if len(ff.followSets) != numSets {
			again = false
		}
	}
}

func (ff *FF) genFollowOf(nt string) *stringset.StringSet {
	// fmt.Printf("genFollowOf(%s)=%s\n", nt, followSets[nt])
	follow := stringset.New()
	for _, r := range ff.g.SyntaxRules {
		for _, a := range r.Alternates {
			bs := a.GetSymbols()
			for _, idx := range stringslice.Find(bs, nt) {
				first := ff.FirstOfString(bs[idx+1:])
				follow.AddSet(first)
				if first.Contain(Empty) {
					// fmt.Printf("  add folow(%s)\n", r.Head.StringValue())
					follow.AddSet(ff.Follow(r.Head.ID()))
				}
			}
		}
	}
	follow.Remove(Empty)
	follow.AddSet(ff.followSets[nt])
	return follow
}

func (ff *FF) initFollowSets() {
	ff.followSets = make(map[string]*stringset.StringSet)
	for _, nt := range ff.g.NonTerminals.Elements() {
		if nt == ff.g.StartSymbol() {
			ff.followSets[nt] = stringset.New("$")
		} else {
			ff.followSets[nt] = stringset.New()
		}
	}
}
