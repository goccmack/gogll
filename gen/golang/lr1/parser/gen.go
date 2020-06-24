/*
This package controls the generation of all parser-related code.
*/
package parser

import (
	"github.com/goccmack/gogll/cfg"
	"github.com/goccmack/gogll/lr1/action"
	"github.com/goccmack/gogll/lr1/basicprod"
	"github.com/goccmack/gogll/lr1/states"
)

func Gen(pkg string, bprods []*basicprod.Production, states *states.States, actions action.Actions) {
	genAction(cfg.BaseDir)
	genActionTable(pkg, cfg.BaseDir, bprods, states, actions)
	genErrors(pkg)
	genGotoTable(cfg.BaseDir, states)
	genParser(pkg, bprods, states)
	genProductionsTable(pkg, bprods, states)

	return
}
