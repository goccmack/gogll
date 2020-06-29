//  Copyright 2020 Marius Ackerman
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
	"text/template"

	"github.com/goccmack/gogll/lr1/states"
	"github.com/goccmack/gogll/symbols"
)

func genGotoRow(state *states.State) string {
	tmpl, err := template.New("parser goto table row").Parse(gotoRowSrc)
	if err != nil {
		panic(err)
	}
	wr := new(bytes.Buffer)
	tmpl.Execute(wr, getGotoRowData(state))
	return wr.String()
}

type gotoRowElement struct {
	NT    string
	State int
}

func getGotoRowData(state *states.State) []gotoRowElement {
	row := make([]gotoRowElement, len(symbols.GetNonTerminalSymbols()))
	for i, nt := range symbols.GetNonTerminalSymbols() {
		row[i].NT = nt
		if nextState := state.Transitions.Transition(nt); nextState == nil {
			row[i].State = -1
		} else {
			row[i].State = nextState.Number
		}
	}
	return row
}

const gotoRowSrc = `{{range $i, $gto := .}}{{$gto.State}}, // {{$gto.NT}}
        {{end}}`
