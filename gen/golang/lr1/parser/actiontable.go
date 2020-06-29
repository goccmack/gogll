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

	"github.com/goccmack/gogll/lr1/action"
	"github.com/goccmack/gogll/lr1/basicprod"
	"github.com/goccmack/gogll/lr1/states"
	"github.com/goccmack/goutil/ioutil"
)

func genActionTable(pkg, outDir string, prods []*basicprod.Production, states *states.States, actions action.Actions) {
	tmpl, err := template.New("parser action table").Parse(actionTableSrc)
	if err != nil {
		panic(err)
	}
	wr := new(bytes.Buffer)
	tmpl.Execute(wr, getActionTableData(pkg, prods, states, actions))
	if err := ioutil.WriteFile(path.Join(outDir, "parser", "actiontable.go"), wr.Bytes()); err != nil {
		panic(err)
	}
}

type actionTableData struct {
	Package string
	Rows    []string
}

func getActionTableData(
	pkg string,
	prods []*basicprod.Production, states *states.States, actions action.Actions,
) (actTab *actionTableData) {
	actTab = &actionTableData{
		Package: pkg,
		Rows:    make([]string, states.Size()),
	}
	for i := range actTab.Rows {
		actTab.Rows[i] = genActionRow(prods, states.List[i], actions[i])
	}
	return
}

const actionTableSrc = `
package parser

import "{{.Package}}/token"

type(
    actionTable [numStates]actionRow
    actionRow struct {
        canRecover bool
        actions map[token.Type]action
    }
)

var actionTab = actionTable{ {{range $i, $r := .Rows}}
	actionRow{ // S{{$i}}
        {{$r}}
	},{{end}}
}

`
