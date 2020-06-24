package parser

import (
	"bytes"
	"fmt"
	"path"
	"strings"
	"text/template"

	"github.com/goccmack/gogll/cfg"
	"github.com/goccmack/gogll/lr1/basicprod"
	"github.com/goccmack/gogll/lr1/states"
	"github.com/goccmack/gogll/symbols"
	"github.com/goccmack/goutil/ioutil"
)

func genProductionsTable(pkg string, prods []*basicprod.Production, states *states.States) {
	fname := path.Join(cfg.BaseDir, "parser", "productionstable.go")
	tmpl, err := template.New("parser productions table").Parse(prodsTabSrc)
	if err != nil {
		panic(err)
	}
	wr := new(bytes.Buffer)
	tmpl.Execute(wr, getProdsTab(pkg, prods, states))
	if err = ioutil.WriteFile(fname, wr.Bytes()); err != nil {
		panic(err)
	}
}

func getProdsTab(pkg string, prods []*basicprod.Production, states *states.States) *prodsTabData {
	data := &prodsTabData{
		Package: pkg,
		ProdTab: make([]prodTabEntry, len(prods)),
	}
	for i, prod := range prods {
		data.ProdTab[i].String = fmt.Sprintf("`%s`", prod.String())
		data.ProdTab[i].Id = prod.Head
		data.ProdTab[i].NTType = int(symbols.GetNTType(prod.Head))
		if len(prod.Body.Symbols) == 0 {
			data.ProdTab[i].NumSymbols = 0
			data.ProdTab[i].ReduceFunc = fmt.Sprintf("nil, nil")
		} else {
			data.ProdTab[i].NumSymbols = len(prod.Body.Symbols)
			data.ProdTab[i].ReduceFunc = fmt.Sprintf("%s%d(%s)",
				prod.Head,
				prod.Alternate,
				getParamIDs(len(prod.Body.Symbols)))
		}
	}

	return data
}

func getParamIDs(n int) string {
	pids := make([]string, n)
	for i := range pids {
		pids[i] = fmt.Sprintf("X[%d]", i)
	}
	return strings.Join(pids, ",")
}

type prodsTabData struct {
	Package string
	ProdTab []prodTabEntry
}

type prodTabEntry struct {
	String     string
	Id         string
	NTType     int
	NumSymbols int
	ReduceFunc string
}

const prodsTabSrc = `
package parser

import(
    "{{.Package}}/ast"
)

type (
	//TODO: change type and variable names to be consistent with other tables
	ProdTab      [numProductions]ProdTabEntry
	ProdTabEntry struct {
		String     string
		Id         string
		NTType     int
		Index int
		NumSymbols int
		ReduceFunc func([]interface{}) (interface{}, error)
	}
)

var productionsTable = ProdTab {
	{{range $i, $entry := .ProdTab}}ProdTabEntry{
		String: {{$entry.String}},
		Id: "{{$entry.Id}}",
		NTType: {{$entry.NTType}},
		Index: {{$i}},
		NumSymbols: {{$entry.NumSymbols}},
		ReduceFunc: func(X []interface{}) (interface{}, error) {
            return ast.{{$entry.ReduceFunc}}
		},
	},
	{{end}}
}
`
