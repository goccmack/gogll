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
Package event implements the events that cause transitions between FSA states.

* Unicode classes: number, letter, upcase, lowcase, space
* Ranges: any, anyof, not
* CharLit
*/
package event

import (
	"fmt"
	"os"
	"sort"
	"unicode"

	"github.com/goccmack/gogll/ast"
	"github.com/goccmack/gogll/lex/item"
	"github.com/goccmack/gogll/util/runeset"
)

// TriState has range: {True, False, Undefined}
type TriState int

const (
	// Undefined is a TriState value
	Undefined TriState = iota
	// False is a TriState value
	False
	// True is a TriState value
	True
)

type eventPair struct {
	a, b ast.LexBase
}

// GetOrdered returns the set of unique transition events for items, ordered
// by the event precidence.
func GetOrdered(items ...*item.Item) (events []ast.LexBase) {
	events = getEvents(items...)
	incompatibleEvents := []eventPair{}
	for i := 0; i < len(events)-1; i++ {
		for j := i + 1; j < len(events); j++ {
			if Subset(events[i], events[j]) == Undefined &&
				Subset(events[j], events[i]) == Undefined {
				incompatibleEvents = append(incompatibleEvents,
					eventPair{events[i], events[j]})
			}
		}
	}
	if len(incompatibleEvents) > 0 {
		fail(items, incompatibleEvents)
	}
	sortEvents(events)
	return
}

// Subset returns True if a is a subset of b, False if a is not a subset of b,
// and Undefined if the subset relationship is not defined between a and b
func Subset(a, b ast.LexBase) TriState {
	switch a1 := a.(type) {
	case *ast.Any:
		_, ok := b.(*ast.Any)
		return toTriState(ok)

	case *ast.AnyOf:
		switch b1 := b.(type) {
		case *ast.Any:
			return True
		case *ast.AnyOf:
			return toTriState(a1.Set.Subset(b1.Set))
		case *ast.CharLiteral:
			return toTriState(a1.Set.Contains(b1.Char()))
		case *ast.Not:
			return toTriState(!a1.Set.Subset(b1.Set))
		case *ast.UnicodeClass:
			return toTriState(unicodeClassContains(b1, a1.Set.Elements()...))
		default:
			panic("Invalid")
		}

	case *ast.CharLiteral:
		switch b1 := b.(type) {
		case *ast.Any:
			return True
		case *ast.AnyOf:
			return toTriState(b1.Set.Contains(a1.Char()))
		case *ast.CharLiteral:
			return toTriState(a1.Char() == b1.Char())
		case *ast.Not:
			return toTriState(!b1.Set.Contains(a1.Char()))
		case *ast.UnicodeClass:
			return toTriState(unicodeClassContains(b1, a1.Char()))
		default:
			panic("Invalid")
		}

	case *ast.Not:
		switch b1 := b.(type) {
		case *ast.Any:
			return True
		case *ast.AnyOf:
			return Undefined
		case *ast.CharLiteral:
			return False
		case *ast.Not:
			return toTriState(b1.Set.Subset(a1.Set))
		case *ast.UnicodeClass:
			return Undefined
		default:
			panic("Invalid")
		}

	case *ast.UnicodeClass:
		switch b1 := b.(type) {
		case *ast.Any:
			return True
		case *ast.AnyOf:
			return Undefined
		case *ast.CharLiteral:
			return False
		case *ast.Not:
			return False
		case *ast.UnicodeClass:
			return toTriState(a1.Type == b1.Type ||
				(b1.Type == ast.Letter && (a1.Type == ast.Lowcase || a1.Type == ast.Upcase)))
		default:
			panic("Invalid")
		}

	default:
		panic("Invalid")
	}
}

func anyOf(r rune, rs *runeset.RuneSet) bool {
	return rs.Contains(r)
}

func contains(events []ast.LexBase, event ast.LexBase) bool {
	for _, ev := range events {
		if ev.Equal(event) {
			return true
		}
	}
	return false
}

func getEvents(items ...*item.Item) (events []ast.LexBase) {
	for _, item := range items {
		event := item.Symbol()
		if event != nil {
			if !contains(events, event.(ast.LexBase)) {
				events = append(events, event.(ast.LexBase))
			}
		}
	}
	return
}

func noneOf(r rune, rs *runeset.RuneSet) bool {
	return !rs.Contains(r)
}

func sortEvents(events []ast.LexBase) {
	sort.Slice(events, func(i, j int) bool {
		switch e1 := events[i].(type) {
		case *ast.Any:
			return false
		case *ast.AnyOf:
			switch e2 := events[j].(type) {
			case *ast.Any:
				return true
			case *ast.AnyOf:
				return e1.Set.Subset(e2.Set)
			case *ast.CharLiteral:
				return e1.Set.Contains(e2.Char())
			case *ast.Not:
				return e1.Set.Subset(e2.Set)
			case *ast.UnicodeClass:
				return unicodeClassContains(e2, e1.Set.Elements()...)
			default:
				panic("Invalid")
			}
		case *ast.CharLiteral:
			switch e2 := events[j].(type) {
			case *ast.Any:
				return true
			case *ast.AnyOf:
				return true
			case *ast.CharLiteral:
				return e1.Char() < e2.Char()
			case *ast.Not:
				return true
			case *ast.UnicodeClass:
				return true
			default:
				panic("Invalid")
			}
		case *ast.Not:
			switch e2 := events[j].(type) {
			case *ast.Any:
				return true
			case *ast.AnyOf:
				return false
			case *ast.CharLiteral:
				return false
			case *ast.Not:
				return e1.Set.Subset(e2.Set)
			case *ast.UnicodeClass:
				return false
			default:
				panic("Invalid")
			}
		case *ast.UnicodeClass:
			switch e2 := events[j].(type) {
			case *ast.Any:
				return true
			case *ast.AnyOf:
				return false
			case *ast.CharLiteral:
				return false
			case *ast.Not:
				return !unicodeClassContains(e1, e2.Set.Elements()...)
			case *ast.UnicodeClass:
				return e2.Type == ast.Letter && (e1.Type == ast.Lowcase || e1.Type == ast.Upcase)
			default:
				panic("Invalid")
			}
		default:
			panic("Invalid")
		}
		panic("Missing return")
	})
}

func toTriState(b bool) TriState {
	if b {
		return True
	}
	return False
}

func unicodeClassContains(c *ast.UnicodeClass, rs ...rune) bool {
	for _, r := range rs {
		ok := false
		switch c.Type {
		case ast.Letter:
			ok = unicode.IsLetter(r)
		case ast.Upcase:
			ok = unicode.IsUpper(r)
		case ast.Lowcase:
			ok = unicode.IsLower(r)
		case ast.Number:
			ok = unicode.IsNumber(r)
		case ast.Space:
			ok = unicode.IsSpace(r)
		default:
			panic("Invalid")
		}
		if !ok {
			return false
		}
	}
	return true
}

func fail(items []*item.Item, incomatibleEvents []eventPair) {
	fmt.Println("Error in lexer events")
	fmt.Println("  Set:")
	for _, item := range items {
		fmt.Println("    ", item)
	}
	fmt.Println("  Incompatible events:")
	for _, ee := range incomatibleEvents {
		fmt.Println("    ", ee.a, " ", ee.b)
	}
	os.Exit(1)
}
