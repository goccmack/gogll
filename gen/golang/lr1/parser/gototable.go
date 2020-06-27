package parser

import (
	"bytes"
	"path"
	"text/template"

	"github.com/goccmack/gogll/lr1/states"
	"github.com/goccmack/gogll/symbols"
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
