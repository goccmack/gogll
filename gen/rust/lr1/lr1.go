/*
Package lr1 generates Rust code for the LR(1) parser
*/
package lr1

import (
	"github.com/goccmack/gogll/gen/rust/lr1/ast"
	"github.com/goccmack/gogll/gen/rust/lr1/parser"
	"github.com/goccmack/gogll/lr1/action"
	"github.com/goccmack/gogll/lr1/basicprod"
	"github.com/goccmack/gogll/lr1/states"
)

func Gen(pkg string, bprods []*basicprod.Production, states *states.States, actions action.Actions) {
	ast.Gen(pkg, bprods)
	parser.Gen(pkg, bprods, states, actions)
}
