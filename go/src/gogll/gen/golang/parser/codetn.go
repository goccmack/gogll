package parser

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
	fmt.Fprintf(buf, "            cU = create(labels.%s, cU, cI, cN)\n", sl.Label())
	fmt.Fprintf(buf, "            L = labels.J_%s\n", sym)
	fmt.Fprintf(buf, "        case labels.%s: // %s\n", sl.Label(), sl.String())
	return buf.String()
}

func codeTNT(sym string, sl gslot.SlotLabel) string {
	buf := new(bytes.Buffer)
	fmt.Fprintf(buf, "            cR = sppf.GetNodeT(`%s`, cI, runeSize)\n", sym)
	fmt.Fprintf(buf, "            cI += runeSize\n")
	fmt.Fprintf(buf, "            next, nextRune, runeSize = decodeRune(input[cI:])\n")
	fmt.Fprintf(buf, "            cN = sppf.GetNode(labels.%s, cN, cR)\n", sl.Label())
	return buf.String()
}
