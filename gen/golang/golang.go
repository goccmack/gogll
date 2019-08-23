package golang

import (
	"gogll/ast"
	"gogll/cfg"
	"gogll/frstflw"
	"gogll/gen/golang/goutil"
	"gogll/gen/golang/parser"
	"gogll/gslot"
	"path"
)

func Gen(g *ast.Grammar, gs *gslot.GSlot, ff *frstflw.FF) {
	goutil.Gen(path.Join(cfg.BaseDir, "goutil"), g)
	parser.Gen(path.Join(cfg.BaseDir, "parser"), g, gs, ff)
}
