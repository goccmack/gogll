package golang

import (
	"gogll/ast"
	"gogll/gen/golang/parser"
	"gogll/gen/golang/sppf"
	"path"
)

func Gen(baseDir string, grammar *ast.Grammar) {
	parserDir := path.Join(baseDir, "parser")
	parser.Gen(parserDir, grammar)
	sppf.Gen(parserDir)
}
