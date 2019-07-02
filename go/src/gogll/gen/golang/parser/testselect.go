package parser

import (
	"bytes"
	"fmt"
	"gogll/ast"
	"gogll/goutil/stringset"
)

func getTestSelectConditions(nt string, str ...string) string {
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
	buf := new(bytes.Buffer)
	for i, c := range conds.Elements() {
		if i < len(conds.Elements())-1 {
			buf.WriteString(c + " ||\n")
		} else {
			buf.WriteString(c)
		}
	}
	return buf.String()
}

func getTestSelectCondition(symName string) string {
	if symName == "$" {
		return "next == Dollar"
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
