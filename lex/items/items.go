/*
Package items computes the lexical item sets, following

Modern Compiler Design. Second Edition.
Grune et al
Springer 2012
Section 2.6
*/
package items

import (
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
		changed = false
		for _, newSet := range sets.Set(i).nextSets() {
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
	return sets
}

// Accept returns the token type accepted by the first reduce item in set
// slits is the set of string literals from the AST
func (set *Set) Accept(slits *stringset.StringSet) string {
	// acceptItems is sorted with string literals first
	acceptItems := set.acceptItems(slits)

	// Check for accepting multiple string literals
	if len(acceptItems) > 1 && slits.Contain(acceptItems[1].Rule.ID()) {
		fmt.Printf("Error in lex item sets: S%d accepts multiple string literals", set.No)
		os.Exit(1)
	}

	if len(acceptItems) > 0 {
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
	events := event.GetOrdered(set.Items()...)
	for _, ev := range events {
		newSet := &Set{}
		for _, item := range set.set {
			if sym := item.Symbol(); sym != nil {
				if event.Subset(ev, sym.(ast.LexBase)) == event.True {
					// if event.Equal(sym.(ast.LexBase)) {
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
	// fmt.Printf("Sets.Contains\n")
	// for _, item := range set.set {
	// 	fmt.Printf("  %s\n", item)
	// }

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
	for _, sl := range g.StringLiterals.ElementsSorted() {
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

func (sets *Sets) add(set *Set) *Sets {
	sets.sets = append(sets.sets, set)
	return sets
}

func stringLitToRule(sl string) *ast.LexRule {
	return &ast.LexRule{ast.StringLitToTokID(sl), stringLitToRegExp(sl)}
}

func stringLitToRegExp(sl string) *ast.RegExp {
	return &ast.RegExp{stringLitToLexSymbols(sl)}
}

func stringLitToLexSymbols(sl string) (symbols []ast.LexSymbol) {
	for _, r := range []rune(sl) {
		symbols = append(symbols, ast.RuneToCharLit(r))
	}
	return
}
