package golang

import (
	"bytes"
	"fmt"
)

func codeAlt(nt string, symbols ...string) string {
	if len(symbols) == 0 {
		return codeTNEmpty
	}
	indent := "			"
	buf := new(bytes.Buffer)
	fmt.Fprintf(buf, codeTN(symbols[0]))
	for i := 1; i < len(symbols); i++ {
		fmt.Fprintf(buf, "%s%s\n", indent, getTestSelectForSymbol(nt, symbols[i]))
		fmt.Fprintf(buf, "%s%s\n", indent, codeTN(symbols[i]))
	}
	fmt.Fprintf(buf, testFollowCode(nt))
	return buf.String()
}

const codeTNEmpty = `cR, cN = getNode(cI), getNode("ϵ", cN, cR)
			pop(cU, cI, cN)`
