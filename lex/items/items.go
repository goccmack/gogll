/*
Copyright 2020 Marius Ackerman

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

  http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

/*
Package items computes the lexical item sets, following

Modern Compiler Design. Second Edition.
Grune et al
Springer 2012
Section 2.6
*/
package items

import (
	"bytes"
	"fmt"
	"os"
	"sort"

	"github.com/goccmack/gogll/ast"
	"github.com/goccmack/gogll/lex/item"
	"github.com/goccmack/gogll/lex/items/event"
	"github.com/goccmack/goutil/stringset"
)

type Event interface {
}

type Set struct {
	No          int
	set         []*item.Item
	Transitions []*Transition
}

type Sets struct {
	sets []*Set
}

type Transition struct {
	Event ast.LexBase
	To    *Set
}

func New(g *ast.GoGLL) *Sets {
	s0 := set0(g)
	s0.No = 0
	sets := new(Sets).add(s0)
	i, changed := 0, true
	for changed || i < sets.Len() {
		// fmt.Printf("item.New: %d sets\n", len(sets.sets))
		changed = false
		for _, newSet := range sets.Set(i).nextSets() {
			// fmt.Printf("  Set %d\n", j)
			if oldSet := sets.GetExisting(newSet); oldSet == nil {
				newSet.No = len(sets.sets)
				sets.add(newSet)
				changed = true
			} else {
				sets.Set(i).changeToSetNo(newSet, oldSet)
			}
		}
		i++
	}
	for _, s := range sets.sets {
		sort.Sort(s)
	}
	// fmt.Println("items.New: done")
	return sets
}

// Accept returns the token type accepted by the first reduce item in set
// slits is the set of string literals from the AST
func (set *Set) Accept(slits *stringset.StringSet) string {
	// acceptItems is sorted with string literals first
	acceptItems := set.acceptItems(slits)

	// Check for accepting multiple string literals
	if len(acceptItems) > 1 && slits.Contain(acceptItems[1].Rule.ID()) {
		fmt.Printf("Error in lex item sets: S%d accepts multiple string literals\n", set.No)
		os.Exit(1)
	}

	if len(acceptItems) > 0 {
		fmt.Printf("items.Accept: %s\n", acceptItems[0].Rule)
		return acceptItems[0].Rule.ID()
	}

	return "Error"
}

// slits is the set of string literals from the AST
func (set *Set) acceptItems(slits *stringset.StringSet) (items []*item.Item) {
	for _, itm := range set.Items() {
		if itm.IsReduce() {
			items = append(items, itm)
		}
	}
	sort.Slice(items, func(i, j int) bool {
		return slits.Contain(items[i].Rule.ID()) &&
			!slits.Clone().Contain(items[j].Rule.ID())
	})
	return
}

func (set *Set) Add(items ...*item.Item) {
	for _, item := range items {
		if !set.Contains(item) {
			set.set = append(set.set, item)
		}
	}
}

func (set *Set) Contains(item *item.Item) bool {
	for _, item1 := range set.set {
		if item.Equal(item1) {
			return true
		}
	}
	return false
}

func (set *Set) Equals(other *Set) bool {
	if len(set.set) != len(other.set) {
		return false
	}
	for _, item := range set.set {
		if !other.Contains(item) {
			return false
		}
	}
	return true
}

func (set *Set) Items() []*item.Item {
	return set.cloneItems()
}

func (set *Set) String() string {
	w := new(bytes.Buffer)
	fmt.Fprintf(w, "S%d:\n", set.No)
	for _, item := range set.Items() {
		fmt.Fprintf(w, "  %s\n", item)
	}
	return w.String()
}

func (set *Set) changeToSetNo(from, to *Set) {
	for _, t := range set.Transitions {
		if t.To == from {
			t.To = to
		}
	}
}

/*
nextSets returns the next set for each possible event transition in set
*/
func (set *Set) nextSets() (sets []*Set) {
	// fmt.Println("items.nextSets")
	// fmt.Println(set)

	events := event.GetOrdered(set.Items()...)
	for _, ev := range events {
		newSet := &Set{}
		for _, item := range set.set {
			if sym := item.Symbol(); sym != nil {
				if event.Subset(ev, sym.(ast.LexBase)) == event.True {
					newSet.Add(item.Next().Emoves()...)
				}
			}
		}
		set.Transitions = append(set.Transitions,
			&Transition{
				Event: ev,
				To:    newSet,
			})
		sets = append(sets, newSet)
	}
	return
}

func (sets *Sets) GetExisting(set *Set) *Set {
	for _, set1 := range sets.sets {
		if set1.Equals(set) {
			return set1
		}
	}
	return nil
}

// Len returns the number of sets in sets
func (sets *Sets) Len() int {
	return len(sets.sets)
}

func (sets *Sets) Set(i int) *Set {
	return sets.sets[i]
}

func (sets *Sets) Sets() []*Set {
	return sets.sets
}

func set0(g *ast.GoGLL) *Set {
	s0 := &Set{}
	for _, rule := range g.LexRules {
		s0.add(item.New(rule).Emoves()...)
	}
	for _, sl := range g.StringLiterals {
		s0.add(item.New(stringLitToRule(sl)))
	}
	return s0
}

func (set *Set) add(items ...*item.Item) *Set {
	for _, item := range items {
		set.set = append(set.set, item)
	}
	return set
}

func (set *Set) cloneItems() []*item.Item {
	items := make([]*item.Item, len(set.set))
	for i, item := range set.set {
		items[i] = item.Clone()
	}
	return items
}

/*** Sort Interface for Set ***/

func (s *Set) Len() int {
	return len(s.set)
}

func (s *Set) Less(i, j int) bool {
	return s.set[i].Rule.ID() < s.set[j].Rule.ID()
}

func (s *Set) Swap(i, j int) {
	s.set[i], s.set[j] = s.set[j], s.set[i]
}

//-------------

func (sets *Sets) add(set *Set) *Sets {
	sets.sets = append(sets.sets, set)
	return sets
}

func stringLitToRule(sl *ast.StringLit) *ast.LexRule {
	return &ast.LexRule{false, ast.StringLitToTokID(sl), stringLitToRegExp(sl)}
}

func stringLitToRegExp(sl *ast.StringLit) *ast.RegExp {
	return &ast.RegExp{stringLitToLexSymbols(sl)}
}

func stringLitToLexSymbols(sl *ast.StringLit) (symbols []ast.LexSymbol) {
	slit := sl.Literal()
	for i := 1; i < len(slit)-1; i++ {
		if slit[i] == '\\' {
			symbols = append(symbols, ast.CharLitFromStringLit(sl, i, true))
			i++
		} else {
			symbols = append(symbols, ast.CharLitFromStringLit(sl, i, false))
		}
	}

	return
}
