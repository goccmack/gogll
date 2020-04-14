package lexer

import (
	"bytes"
	"path/filepath"
	"text/template"

	"github.com/goccmack/gogll/ast"
	"github.com/goccmack/goutil/ioutil"
)

type Data struct {
	Package string
}

func Gen(lexDir string, g *ast.GoGLL) {
	tmpl, err := template.New("lexer").Parse(tmplSrc)
	if err != nil {
		panic(err)
	}
	buf := new(bytes.Buffer)
	if err = tmpl.Execute(buf, getData(g)); err != nil {
		panic(err)
	}
	if err = ioutil.WriteFile(filepath.Join(lexDir, "lexer.go"), buf.Bytes()); err != nil {
		panic(err)
	}
}

func getData(g *ast.GoGLL) *Data {
	return &Data{
		Package: g.Package.GetString(),
	}
}

const tmplSrc = `
// Package lexer is fenerated by GoGLL. Do not edit.
package lexer

import(
    "io/ioutil"
    "strings"

    "github.com/goccmack/goutil/md"

    "{{.Package}}/token"
)

type Lexer struct {
    I []rune
    Tokens []*token.Token
}

func NewFile(fname string) *Lexer {
    if strings.HasSuffix(fname, ".md") {
        src, err := md.GetSource(fname)
        if err != nil {
            panic(err)
        }
        return New([]rune(src))
    }
    buf, err := ioutil.ReadFile(fname)
    if err != nil {
        panic(err)
    }
    return New([]rune(string(buf)))
}

func New(input []rune) *Lexer {
    panic("implement")
}

// GetLineColumn returns the line and column of rune[i] in the input
func (l *Lexer) GetLineColumn(i int) (line, col int) {
	line, col = 1, 1
	for j := 0; j < i; j++ {
		switch l.I[j] {
		case '\n':
			line++
			col = 1
		case '\t':
			col += 4
		default:
			col++
		}
	}
	return
}
`
