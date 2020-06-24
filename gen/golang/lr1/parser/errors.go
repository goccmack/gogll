package parser

import (
	"bytes"
	"path"
	"text/template"

	"github.com/goccmack/gogll/cfg"
	"github.com/goccmack/goutil/ioutil"
)

func genErrors(pkg string) {
	tmpl, err := template.New("parser errors").Parse(errorsSrc)
	if err != nil {
		panic(err)
	}
	wr := new(bytes.Buffer)
	tmpl.Execute(wr, path.Join(pkg, "token"))
	if err := ioutil.WriteFile(path.Join(cfg.BaseDir, "errors", "errors.go"), wr.Bytes()); err != nil {
		panic(err)
	}
}

const errorsSrc = `
package errors

import(
	"bytes"
	"fmt"
	"{{.}}"
)

type ErrorSymbol interface {
}

type Error struct {
	Err            error
	ErrorToken     *token.Token
	ErrorSymbols   []ErrorSymbol
	ExpectedTokens []string
}

func (E *Error) String() string {
	w := new(bytes.Buffer)
	fmt.Fprintf(w, "Error")
	if E.Err != nil {
		fmt.Fprintf(w, " %s\n", E.Err)
	} else {
		fmt.Fprintf(w, "\n")
	}
	fmt.Fprintf(w, "Token: type=%d, lit=%s\n", E.ErrorToken.Type, E.ErrorToken.Literal())
	ln, col := E.ErrorToken.GetLineColumn()
	fmt.Fprintf(w, "Pos: offset=%d, line=%d, column=%d\n", E.ErrorToken.Lext(), ln, col)
	fmt.Fprintf(w, "Expected one of: ")
	for _, sym := range E.ExpectedTokens {
		fmt.Fprintf(w, "%s ", sym)
	}
	fmt.Fprintf(w, "ErrorSymbol:\n")
	for _, sym := range E.ErrorSymbols {
		fmt.Fprintf(w, "%v\n", sym)
	}
	return w.String()
}
`
