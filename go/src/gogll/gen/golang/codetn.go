package golang

import (
	"bytes"
	"fmt"
	"gogll/ast"
	"gogll/gslot"
)

func codeTN(sym string, sl gslot.SlotLabel) string {
	if ast.IsTerminal(sym) {
		return codeTNT(sym, sl)
	}
	return codeTNNT(sym, sl)
}

func codeTNNT(sym string, sl gslot.SlotLabel) string {
	buf := new(bytes.Buffer)
	fmt.Fprintf(buf, "            cU = stack.Create(%s, cU, cI, cN)\n", sl.Label())
	fmt.Fprintf(buf, "        case %s: // %s\n", sl.Label(), sl.String())
	return buf.String()
}

func codeTNT(sym string, sl gslot.SlotLabel) string {
	buf := new(bytes.Buffer)
	fmt.Fprintf(buf, "            cR := forest.GetNodeT(\"%s\", cI)\n", sym)
	fmt.Fprintf(buf, "            cI += runeSize\n")
	fmt.Fprintf(buf, "            cN = forest.GetNode(%s, cN, cR)\n", sl.Label())
	return buf.String()
}
