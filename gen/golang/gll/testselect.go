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

package gll

import (
	"bytes"
	"fmt"
	"os"
	"sort"
	"text/template"

	"github.com/goccmack/gogll/v3/frstflw"
	"github.com/goccmack/gogll/v3/gen/golang/utils"
	"github.com/goccmack/gogll/v3/gslot"
	"github.com/goccmack/gogll/v3/symbols"
)

func (g *gen) genTestSelect() string {
	tmpl, err := template.New("Test Select").Parse(testSelectTmpl)
	if err != nil {
		panic(err)
	}
	buf, data := new(bytes.Buffer), g.getTestSelectData()
	if err = tmpl.Execute(buf, data); err != nil {
		panic(err)
	}
	return buf.String()
}

type TestSelectData struct {
	TestSelect []*TSData
	Follow     []*TSData
}

type Symbol struct {
	TokType string
	Label   string
}

type TSData struct {
	Label   string
	Symbols Symbols
	// Conditions []*Condition
}

type Symbols []*Symbol

// type Condition struct {
// 	Cond string
// 	Last bool
// }

func (g *gen) getTestSelectData() *TestSelectData {
	return &TestSelectData{
		TestSelect: g.getTSData(),
		Follow:     g.getFollowData(),
	}
}

func (g *gen) getFollowData() (data []*TSData) {
	for _, nt := range g.g.NonTerminals.ElementsSorted() {
		data = append(data, g.getFollowDataForNT(nt))
	}
	return
}

func (g *gen) getFollowDataForNT(nt string) *TSData {
	d := &TSData{
		Label:   nt,
		Symbols: g.getFollowConditions(nt),
	}
	// fmt.Printf("getFollowDataForNT(%s): %d\n", nt, len(d.Conditions))
	return d
}

func (g *gen) getTSData() (data []*TSData) {
	for _, s := range g.gs.Slots() {
		data = append(data, g.getSlotTSData(s))
	}
	return
}

func (g *gen) getSlotTSData(l gslot.Label) *TSData {
	data := &TSData{
		Label:   l.String(),
		Symbols: g.getFirst(l),
		// Conditions: g.getSlotTSConditions(s),
	}
	checkDuplicateSymbols(data)
	return data
}

func checkDuplicateSymbols(data *TSData) {
	for i := 0; i < len(data.Symbols)-1; i++ {
		for j := i + 1; j < len(data.Symbols); j++ {
			if data.Symbols[i].Label == data.Symbols[j].Label {
				fmt.Printf("Duplicate symbol %s in %s\n",
					data.Symbols[i], data.Label)
				panic("cleanup")
			}
		}
	}
}

func (g *gen) getFirst(l gslot.Label) (tokens Symbols) {
	// fmt.Printf("testSelect.getFirst(%s)\n", l)

	ss := l.Symbols()[l.Pos:]
	frst := g.ff.FirstOfString(ss.Strings())
	firstSymbols := frst.Elements()
	sort.Slice(
		firstSymbols,
		func(i, j int) bool { return firstSymbols[i] < firstSymbols[j] })

	for _, sym := range firstSymbols {
		if sym != frstflw.Empty {
			tokens = append(tokens,
				&Symbol{
					TokType: symbols.TerminalLiteralToType(sym).TypeString(),
					Label:   utils.Escape(sym),
				})
		}
	}
	if frst.Contain(frstflw.Empty) {
		for _, s := range g.getFollowConditions(l.Head) {
			if !tokens.contains(s.Label) {
				tokens = append(tokens, s)
			}
		}
	}

	// fmt.Printf("  first:\n")
	// for _, t := range tokens {
	// 	fmt.Println("    ", t.Label)
	// }

	return
}

func (g *gen) getFollowConditions(nt string) (tokens []*Symbol) {
	// fmt.Printf("testselect.getFollowConditions(%s)\n", nt)
	flw := g.ff.Follow(nt)
	if flw.Len() == 0 {
		fmt.Printf("Production %s has empty follow set. It is never called\n", nt)
		os.Exit(1)
	}
	for _, sym := range flw.ElementsSorted() {
		// fmt.Printf("getFollowConditions: %s\n", sym)
		tokens = append(tokens,
			&Symbol{
				TokType: symbols.TerminalLiteralToType(sym).TypeString(),
				Label:   utils.Escape(sym),
			})
	}
	return
}

func (ss Symbols) contains(label string) bool {
	for _, s := range ss {
		if s.Label == label {
			return true
		}
	}
	return false
}

const testSelectTmpl = `
var first = []map[token.Type]string { {{range $ts := .TestSelect}}
	// {{$ts.Label}}
	{ {{range $sym := $ts.Symbols}}
		token.{{$sym.TokType}}:"{{$sym.Label}}",{{end}}
	},{{end}}
}

var followSets = []map[token.Type]string { {{range $flw := .Follow}}
	// {{$flw.Label}}
	{ {{range $sym := $flw.Symbols}}
		token.{{$sym.TokType}}:"{{$sym.Label}}",{{end}}
	},{{end}}
} 
`
