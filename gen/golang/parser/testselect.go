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

package parser

import (
	"bytes"
	"fmt"
	"os"
	"sort"
	"text/template"

	"github.com/goccmack/gogll/frstflw"
	"github.com/goccmack/gogll/gen/golang/token"
	"github.com/goccmack/gogll/gslot"
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
	Symbols []*Symbol
	// Conditions []*Condition
}

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
	return &TSData{
		Label:   l.String(),
		Symbols: g.getFirst(l),
		// Conditions: g.getSlotTSConditions(s),
	}
}

func (g *gen) getFirst(l gslot.Label) (tokens []*Symbol) {
	// fmt.Printf("testSelect.getFirst(%s)\n", l)
	ss := l.Symbols()[l.Pos:]
	frst := g.ff.FirstOfString(ss.Strings())
	firstSymbols := frst.Elements()
	sort.Slice(
		firstSymbols,
		func(i, j int) bool { return firstSymbols[i] < firstSymbols[j] })
	tokMap := token.GetTokenMap(g.g)
	// fmt.Printf("  first: %s\n", frst)
	for _, sym := range firstSymbols {
		if sym != frstflw.Empty {
			tokens = append(tokens,
				&Symbol{
					TokType: tokMap[sym],
					Label:   sym,
				})
		}
	}
	if frst.Contain(frstflw.Empty) {
		tokens = append(tokens, g.getFollowConditions(l.Head)...)
	}
	return
}

func (g *gen) getFollowConditions(nt string) (tokens []*Symbol) {
	// fmt.Printf("testselect.getFollowConditions(%s)\n", nt)
	flw := g.ff.Follow(nt)
	if flw.Len() == 0 {
		fmt.Printf("Production %s has empty follow set. It is never called\n", nt)
		os.Exit(1)
	}
	tokMap := token.GetTokenMap(g.g)
	for _, sym := range flw.Elements() {
		tokens = append(tokens,
			&Symbol{
				TokType: tokMap[sym],
				Label:   sym,
			})
	}
	return
}

const testSelectTmpl = `
var first = []map[token.Type]string { {{range $ts := .TestSelect}}
    // {{$ts.Label}}
    map[token.Type]string{ {{range $sym := $ts.Symbols}}
        token.{{$sym.TokType}}:"{{$sym.Label}}",{{end}}
    },{{end}}
}

var followSets = []map[token.Type]string { {{range $flw := .Follow}}
    // {{$flw.Label}}
    map[token.Type]string{ {{range $sym := $flw.Symbols}}
        token.{{$sym.TokType}}:"{{$sym.Label}}",{{end}}
    },{{end}}
} 
`

/*
var testSelect = []func()bool { {{range $i, $ts := .TestSelect}}
    // slot.{{$ts.Label}}
    func()bool{
        return {{range $i, $c := $ts.Conditions}}{{$c.Cond}} {{if not $c.Last}}||{{end}}
    {{end}} },
{{end}} }

{{range $i, $flw := .Follow}}func follow{{$flw.Label}} () bool {
    return {{range $i, $c := $flw.Conditions}}{{$c.Cond}} {{if not $c.Last}}||{{end}}
{{end}} }
*/
