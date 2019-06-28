package golang

import (
	"bytes"
	"fmt"
	"gogll/ast"
	"gogll/gslot"
)

func codeAlt(rule *ast.Rule, altI int) string {
	symbols := rule.Alternates[altI].Symbols()
	if symbols[0] == ast.Empty {
		return codeTNEmpty(gslot.SlotLabel{rule.Head.Value(), altI, 1})
	}
	indent := "			"
	buf := new(bytes.Buffer)
	fmt.Fprintf(buf, codeTN(symbols[0], gslot.SlotLabel{rule.Head.Value(), altI, 1}))
	for i := 1; i < len(symbols); i++ {
		fmt.Fprintf(buf, "%s%s\n", indent, getTestSelectForSymbol(rule.Head.Value(), symbols[i]))
		sl := gslot.SlotLabel{rule.Head.Value(), altI, i + 1}
		fmt.Fprintf(buf, "%s%s\n", indent, codeTN(symbols[i], sl))
	}
	fmt.Fprintf(buf, testFollowCode(rule.Head.Value()))
	return buf.String()
}

func codeTNEmpty(sl gslot.SlotLabel) string {
	code := fmt.Sprintf(`cR := forest.GetNodeE(cI)
			cN = forest.GetNode(%s, cN, cR)
			stack.Pop(cU, cI, cN)`, sl.Label())
	return code
}
