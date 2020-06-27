package parser

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/goccmack/gogll/lr1/action"
	"github.com/goccmack/gogll/lr1/basicprod"
	"github.com/goccmack/gogll/lr1/states"
	"github.com/goccmack/gogll/symbols"
)

func genActionRow(prods []*basicprod.Production, state *states.State, actions map[string]action.Action) string {
	wr := new(bytes.Buffer)
	tmpl, err := template.New("parser action row").Parse(actionRowSrc)
	if err != nil {
		panic(err)
	}
	tmpl.Execute(wr, getActionRowData(prods, state, actions))
	return wr.String()
}

type actRow struct {
	CanRecover bool
	Actions    []*Action
}

type Action struct {
	Token  string
	Action string
}

func getActionRowData(prods []*basicprod.Production, state *states.State, actions map[string]action.Action) (data *actRow) {
	data = &actRow{
		CanRecover: state.CanRecover(),
		Actions:    []*Action{},
	}
	for _, sym := range symbols.GetTerminals() {
		if actions[sym.Literal()] != nil {
			var actStr string
			switch act := actions[sym.Literal()].(type) {
			case action.Accept:
				actStr = fmt.Sprintf("accept(true),\t\t/* %s */", sym)
			case action.Reduce:
				actStr = fmt.Sprintf("reduce(%d),\t\t/* %s, reduce: %s */", int(act), sym, prods[int(act)].Head)
			case action.Shift:
				actStr = fmt.Sprintf("shift(%d),\t\t/* %s */", int(act), sym)
			default:
				panic(fmt.Sprintf("Unknown action type: %T", act))
			}
			data.Actions = append(data.Actions,
				&Action{
					Token:  sym.GoString(),
					Action: actStr,
				})
		}
	}
	return
}

const actionRowSrc = `canRecover: {{printf "%t" .CanRecover}},
		actions: map[token.Type]action{ {{range $a := .Actions}}
			token.{{$a.Token}}:{{$a.Action}}{{end}}
        },
`
