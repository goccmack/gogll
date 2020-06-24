package codegen

import (
	"github.com/goccmack/gogll/lr1/action"
	"github.com/goccmack/gogll/lr1/basicprod"
	"github.com/goccmack/gogll/lr1/codegen/golang"
	"github.com/goccmack/gogll/lr1/states"
)

func Gen(pkg string, prods []*basicprod.Production, states *states.States, actions action.Actions) {
	golang.Gen(pkg, prods, states, actions)
}
