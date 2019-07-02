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
		return fmt.Sprintf("not(nextRune, %s)", string(s.Token.Lit))
	case *ast.Space:
		return "space(nextRune)"
	case *ast.CharLiteral:
		if s.Rune == '"' {
			return "next == \"\\\""
		}
		return fmt.Sprintf("next == \"%c\"", s.Rune)
	case *ast.UpCase:
		return "upcase(nextRune)"
	case *ast.LowCase:
		return "lowcase(nextRune)"
	case *ast.Letter:
		return "letter(nextRune)"
	case *ast.Number:
		return "number(nextRune)"
	case *ast.StringChar:
		if s.Value() == "\"" {
			return "next == \"\\\""
		}
		return fmt.Sprintf("next == \"%s\"", s.Value())
	}
	panic(fmt.Sprintf("Invalid symbol %s", symName))
}
