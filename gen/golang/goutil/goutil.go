package goutil

import (
	"gogll/ast"
	"gogll/gen/golang/goutil/bsr"
	"gogll/gen/golang/goutil/md"
	"gogll/gen/golang/goutil/stringset"
	"path/filepath"
)

func Gen(utilDir string, g *ast.Grammar) {
	bsr.Gen(filepath.Join(utilDir, "bsr", "bsr.go"), g.Package)
	stringset.Gen(filepath.Join(utilDir, "stringset", "stringset.go"))
	md.Gen(filepath.Join(utilDir, "md", "md.go"))
}
