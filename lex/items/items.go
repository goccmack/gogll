/*
Package items computes the lexical item sets, following

Modern Compiler Design. Second Edition.
Grune et al
Springer 2012
Section 2.6
*/
package items

import (
	"github.com/goccmack/gogll/ast"
	"github.com/goccmack/gogll/lex/item"
	"github.com/goccmack/gogll/lex/items/event"
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
	s0 := set0(g.LexRules)
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
	for _, event := range events {
		newSet := &Set{}
		for _, item := range set.set {
			if sym := item.Symbol(); sym != nil {
				if event.Equal(sym.(ast.LexBase)) {
					newSet.Add(item.Next().Emoves()...)
				}
			}
		}
		set.Transitions = append(set.Transitions,
			&Transition{
				Event: event,
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

func set0(rules []*ast.LexRule) *Set {
	s0 := &Set{}
	for _, rule := range rules {
		s0.add(item.New(rule).Emoves()...)
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
