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
