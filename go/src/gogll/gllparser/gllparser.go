package gllparser

import (
	"go/token"
	"gogll2/descriptors"
	"gogll2/gss"
	"gogll2/input"
	"gogll2/popped"
	"gogll2/sppf"
)

func Parse(in *input.Input) {
	u0 := gss.Node{"L0", 0}
	cI := 0
	cU := u0
	cN := sppf.Dummy
	U := descriptors.Empty()
	R := descriptors.Empty()
	P := popped.Empty()
	label := "J_Symbol"

	for {
		switch label {
		case "L0":
		case "J_Symbol":
		}
	}
}

func testSelect(a *token.Token)
