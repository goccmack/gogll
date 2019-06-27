package golang

import "gogll/ast"

func codeTN(sym string) string {
	if ast.IsTerminal(sym) {
		return codeTNT(sym)
	}
	return codeTNNT(sym)
}

func codeTNNT(sym string) string {
	panic("implement me")
}

func codeTNT(sym string) string {
	panic("implement me")
}
