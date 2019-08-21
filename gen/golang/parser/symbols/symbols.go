/*
Package symbols generates parser/symbols.go
*/
package symbols

import (
	"bytes"
	"gogll/ast"
	"gogll/goutil/ioutil"
	"text/template"
)

func Gen(fname string) {
	tmpl, err := template.New("symbols").Parse(src)
	if err != nil {
		panic(err)
	}
	buf := new(bytes.Buffer)
	if err := tmpl.Execute(buf, ast.GetNonTerminals()); err != nil {
		panic(err)
	}
	if err := ioutil.WriteFile(fname, buf.Bytes()); err != nil {
		panic(err)
	}
}

const src = `
package symbols

func IsNonTerminal(symbol string) bool {
	return nonTerminals[symbol]
}

func IsTerminal(symbol string) bool {
	return !nonTerminals[symbol]
}

var nonTerminals = map[string]bool{ {{range $i, $sym := .}}
	"{{$sym}}":true,{{end}}
}
`
