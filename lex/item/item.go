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
Package item implements a lexical dotted item
*/
package item

import (
	"bytes"
	"fmt"

	"github.com/goccmack/gogll/ast"
	"github.com/goccmack/gogll/lex/item/pos"
)

// Item contains the dotted item for one lex rule
type Item struct {
	Rule *ast.LexRule
	Pos  *pos.Pos
}

func From(rule *ast.LexRule, pos *pos.Pos) *Item {
	return &Item{
		Rule: rule,
		Pos:  pos,
	}
}

func New(rule *ast.LexRule) *Item {
	return &Item{
		Rule: rule,
		Pos:  pos.New(),
	}
}

// afterCurrentItem returns the item with position after the current bracket, e.g.:
// X -> 𝜶 (β•| 𝞬) δ => X -> 𝜶 (β | 𝞬) •δ
// X -> 𝜶 (β | 𝞬•) δ => X -> 𝜶 (β | 𝞬) •δ
func (i *Item) afterCurrentBracket() *Item {
	newItem := i.beforeCurrentBracket()
	newItem.Pos.Inc()
	return newItem
}

// beforeCurrentItem returns the item with position before the current bracket, e.g.:
// X -> 𝜶 (β•| 𝞬) => X -> 𝜶 •(β | 𝞬)
// X -> 𝜶 (β | 𝞬•) => X -> 𝜶 •(β | 𝞬)
func (i *Item) beforeCurrentBracket() *Item {
	return &Item{
		Rule: i.Rule,
		Pos:  i.Pos.Clone().Pop(2),
	}
}

// isEndOfBracket returns true iff: the next symbol is a LexBracket
func (i *Item) atBeforeBracket() bool {
	_, ok := i.Symbol().(*ast.LexBracket)
	return ok
}

// isEndOfBracket returns true iff:
// * The current regexp is and alternate of a bracket expression;
// * The position after the last symbol of the current alternate of the current bracket.
func (i *Item) atEndOfBracket() bool {
	// fmt.Println("Item.atEndOfBracket: ", i, i.Pos)
	// fmt.Println("  top=", i.Pos.Top(), " len=", len(i.GetRegExp().Symbols))

	return i.Pos.Len() > 1 && i.Pos.Top() >= len(i.GetRegExp().Symbols)
}

func (i *Item) getCurrentBracket() *ast.LexBracket {
	item := &Item{
		Rule: i.Rule,
		Pos:  i.Pos.Clone().Pop(2),
	}
	return item.Symbol().(*ast.LexBracket)
}

func (i *Item) Clone() *Item {
	return &Item{
		Rule: i.Rule,
		Pos:  i.Pos.Clone(),
	}
}

func (i *Item) Emoves() []*Item {
	// fmt.Println("Emoves:", i)

	after := []*Item{i}
	for changed := true; changed; {
		// fmt.Println(" Again")
		before := after
		after = []*Item{}
		changed = false
		for _, item := range before {
			// fmt.Println("  ", item)

			switch {
			case item.atBeforeBracket():
				for j := range item.Symbol().(*ast.LexBracket).Alternates {
					after = append(after, &Item{
						Rule: item.Rule,
						Pos:  item.Pos.Clone().Push(j).Push(0),
					})
				}
				changed = true
			case item.atEndOfBracket():
				switch item.getCurrentBracket().Type {
				case ast.LexZeroOrMore, ast.LexOneOrMore:
					after = append(after, item.beforeCurrentBracket())
				}
				after = append(after, item.afterCurrentBracket())
				changed = true
			default:
				after = append(after, item)
			}
		}
	}
	return after
}

// func (i *Item) Emoves() []*Item {
// 	after := []*Item{i}
// 	for changed := true; changed; {
// 		before := after
// 		after = []*Item{}
// 		changed = false
// 		for _, item := range before {
// 			if brkt, ok := item.Symbol().(*ast.LexBracket); ok {
// 				for j := range brkt.Alternates {
// 					after = append(after, &Item{
// 						Rule: item.Rule,
// 						Pos:  item.Pos.Clone().Push(j).Push(0),
// 					})
// 				}
// 				changed = true
// 			} else {
// 				after = append(after, item)
// 			}
// 		}
// 	}
// 	return after
// }

// Equal is true if i and other have the same rule and position.
func (i *Item) Equal(other *Item) bool {
	return i.Rule.ID() == other.Rule.ID() &&
		i.Pos.Equal(other.Pos)
}

// Next returns the next item after i. If i is a reduce item Next returns nil
func (i *Item) Next() *Item {
	if i.IsReduce() {
		return nil
	}
	next := From(i.Rule, i.Pos.Clone().Inc())
	// for next.Pos.Len() > 1 && next.Pos.Top() >= len(next.GetRegExp().Symbols) {
	// 	next.Pos.Pop(2)
	// 	next.Pos.Inc()
	// }
	return next
}

// IsReduce returns true if the position of i is after the last symbol of i
func (i *Item) IsReduce() bool {
	return i.Pos.Len() == 1 && i.Pos.Top() >= len(i.Rule.RegExp.Symbols)
}

// Symbol returns the base (not bracket) symbol after the item place marker.
// If i is a reduce item Symbol returns nil
func (i *Item) Symbol() ast.LexSymbol {
	// fmt.Printf("Item.Symbol: %s %s\n", i.Rule, i.Pos)

	re := i.GetRegExp()

	// fmt.Printf("  %s\n", re)

	if i.Pos.Top() >= len(re.Symbols) {
		return nil
	}
	return re.Symbols[i.Pos.Top()]
}

func (i *Item) String() string {
	// if i.IsReduce() {
	// 	return fmt.Sprintf("%s : %s •",
	// 		i.Rule.ID(), stringRegExp(i.Rule.RegExp, i.Pos, 0))

	// }
	str := fmt.Sprintf("%s : %s %s",
		i.Rule.ID(), i.stringRegExp(i.Rule.RegExp, pos.New()), i.Pos)
	return str
}

// GetRegExp returns the RegExp containing the current symbol in i. If pos.Len() > 1
// this will be an alternate of a LexBracket.
func (i *Item) GetRegExp() *ast.RegExp {
	// fmt.Printf("Item.GetRegExp: %s, %s\n", i.Rule, i.Pos)

	re := i.Rule.RegExp
	for j := 0; j < i.Pos.Len()-1; {

		// fmt.Printf("  re: %s pos %d\n", re.Symbols, i.Pos.Peek(j))

		sym := re.Symbols[i.Pos.Peek(j)]
		j++
		if brkt, ok := sym.(*ast.LexBracket); ok {
			re = brkt.Alternates[i.Pos.Peek(j)]
			j++
		}
	}
	return re
}

func (i *Item) stringRegExp(regExp *ast.RegExp, pos *pos.Pos) string {
	// fmt.Printf("item.stringRegExp %s pos=%s\n", regExp, pos)
	w := new(bytes.Buffer)
	for _, symbol := range regExp.Symbols {
		if i.Pos.Equal(pos) {
			fmt.Fprintf(w, "•")
		}
		fmt.Fprintf(w, i.stringSymbol(symbol, pos))
		pos.Inc()
	}
	if i.Pos.Equal(pos) {
		fmt.Fprintf(w, "•")
	}
	return w.String()
}

func (i *Item) stringSymbol(symbol ast.LexSymbol, pos *pos.Pos) string {
	// fmt.Printf("item.stringSymbol %s pos=%s\n", symbol, pos)
	w := new(bytes.Buffer)
	if brkt, ok := symbol.(*ast.LexBracket); ok {
		fmt.Fprintf(w, brkt.LeftBracket())
		for j, alt := range brkt.Alternates {
			if j > 0 {
				fmt.Fprintf(w, "| ")
			}
			fmt.Fprint(w, i.stringRegExp(alt, pos.Clone().Push(j).Push(0)))
		}
		fmt.Fprintf(w, brkt.RightBracket())
	} else {
		fmt.Fprintf(w, "%s ", symbol)
	}
	return w.String()
}

// func stringRegExp(regExp *ast.RegExp, nextSym ast.LexSymbol) string {
// 	w := new(bytes.Buffer)
// 	for _, symbol := range regExp.Symbols {
// 		if symbol == nextSym {
// 			fmt.Fprintf(w, "•")
// 		}
// 		fmt.Fprintf(w, stringSymbol(symbol, nextSym))
// 	}
// 	return w.String()
// }

// func stringSymbol(symbol, next ast.LexSymbol) string {
// 	w := new(bytes.Buffer)
// 	if brkt, ok := symbol.(*ast.LexBracket); ok {
// 		fmt.Fprintf(w, brkt.LeftBracket())
// 		for i, alt := range brkt.Alternates {
// 			if i > 0 {
// 				fmt.Fprintf(w, "| ")
// 			}
// 			fmt.Fprint(w, stringRegExp(alt, next))
// 		}
// 		fmt.Fprintf(w, brkt.RightBracket())
// 	} else {
// 		fmt.Fprintf(w, "%s ", symbol)
// 	}
// 	return w.String()
// }
