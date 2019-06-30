package parser

import (
	"bytes"
	"fmt"
	"gogll/ast"
	"gogll/goutil/stringset"
	"text/template"
)

func getTestSelectForAlternate(nt string, i int, str ...string) string {
	data := getTestSelectDataForAlternate(nt, i, str...)
	tmpl, err := template.New("testSelect").Parse(testSelectForAltTemplate)
	if err != nil {
		failError(err)
	}
	buf := new(bytes.Buffer)
	if err := tmpl.Execute(buf, data); err != nil {
		failError(err)
	}
	return buf.String()
}

func getTestSelectForSymbol(nt string, sym string) string {
	data := getTestSelectDataForAlternate(nt, 0, sym)
	tmpl, err := template.New("testSelect for Symbol").Parse(testSelectForSymbolTemplate)
	if err != nil {
		failError(err)
	}
	buf := new(bytes.Buffer)
	if err := tmpl.Execute(buf, data); err != nil {
		failError(err)
	}
	return buf.String()
}

type testSelectData struct {
	Conditions []string
	Label      string
}

func getTestSelectDataForAlternate(nt string, i int, str ...string) *testSelectData {
	data := &testSelectData{
		Conditions: getTestSelectConditions(nt, i, str),
		Label:      getAlternateLabel(nt, i),
	}
	return data
}

func getTestSelectConditions(nt string, i int, str []string) []string {
	conds := stringset.New()
	empty := str[0] == ast.Empty
	for _, sym := range g.FirstOfString(str).Elements() {
		if sym == ast.Empty {
			empty = true
		} else {
			conds.Add(getTestSelectCondition(sym))
		}
	}
	if empty {
		for _, sym := range g.Follow(nt).Elements() {
			conds.Add(getTestSelectCondition(sym))
		}
	}
	tsConds := make([]string, 0, conds.Len())
	for i, c := range conds.Elements() {
		if i < len(conds.Elements())-1 {
			tsConds = append(tsConds, c+" ||")
		} else {
			tsConds = append(tsConds, c+" {")
		}
	}
	return tsConds
}

func getTestSelectCondition(symName string) string {
	if symName == "$" {
		return "next == \"$\""
	}
	sym := ast.GetSymbol(symName)
	switch s := sym.(type) {
	case *ast.Head:
		panic(fmt.Sprintf("first returns only terminals but got %s", symName))
	case *ast.ID:
		panic(fmt.Sprintf("first returns only terminals but got %s", symName))
	case *ast.AnyChar:
		return "true"
	case *ast.NotString:
		return fmt.Sprintf("not(next, %s)", string(s.Token.Lit))
	case *ast.Space:
		return "unicode.IsSpace(next)"
	case *ast.CharLiteral:
		if s.Rune == '"' {
			return "next == \"\\\""
		}
		return fmt.Sprintf("next == \"%c\"", s.Rune)
	case *ast.UpCase:
		return "unicode.IsUpper(next)"
	case *ast.LowCase:
		return "unicode.IsLower(next)"
	case *ast.Letter:
		return "unicode.IsLetter(next)"
	case *ast.Number:
		return "unicode.IsNumber(next)"
	case *ast.StringChar:
		if s.Value() == "\"" {
			return "next == \"\\\""
		}
		return fmt.Sprintf("next == \"%s\"", s.Value())
	}
	panic(fmt.Sprintf("Invalid symbol %s", symName))
}

const testSelectForAltTemplate = `if {{range $i, $cond := .Conditions}}{{$cond}}
			{{end}}	add(labels.{{.Label}}, cU, cI, dummy)
			}`

const testSelectForSymbolTemplate = `if {{range $i, $cond := .Conditions}}{{$cond}}
			{{end}}	L = labels.L0 
			}`
