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
	"path"
	"text/template"

	"github.com/goccmack/gogll/v3/lr1/states"
	"github.com/goccmack/gogll/v3/symbols"
	"github.com/goccmack/goutil/ioutil"
)

func genGotoTable(outDir string, states *states.States) {
	tmpl, err := template.New("parser goto table").Parse(gotoTableSrc)
	if err != nil {
		panic(err)
	}
	wr := new(bytes.Buffer)
	if err := tmpl.Execute(wr, getGotoTableData(states)); err != nil {
		panic(err)
	}
	if err := ioutil.WriteFile(path.Join(outDir, "parser", "gototable.go"), wr.Bytes()); err != nil {
		panic(err)
	}
}

func getGotoTableData(states *states.States) *gotoTableData {
	data := &gotoTableData{
		NumNTSymbols: len(symbols.GetNonTerminalSymbols()),
		Rows:         make([]string, len(states.List)),
	}
	for i, state := range states.List {
		data.Rows[i] = genGotoRow(state)
	}
	return data
}

type gotoTableData struct {
	NumNTSymbols int
	Rows         []string
}

const gotoTableSrc = `
/*
*/
package parser

const numNTSymbols = {{.NumNTSymbols}}
type(
	gotoTable [numStates]gotoRow
	gotoRow	[numNTSymbols] int
)

var gotoTab = gotoTable{
	{{range $i, $r := .Rows}}gotoRow{ // S{{$i}}
		{{$r}}
	},
	{{end}}
}
`
