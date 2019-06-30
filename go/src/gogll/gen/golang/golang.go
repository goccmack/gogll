package golang

import (
	"gogll/ast"
	"gogll/gen/golang/parser"
	"gogll/gen/golang/sppf"
	"os"
	"path"
)

func Gen(baseDir string, grammar *ast.Grammar) {
	parserDir := path.Join(baseDir, "parser")
	if err := os.MkdirAll(parserDir, 0731); err != nil {
		panic(err)
	}
	parser.Gen(parserDir, grammar)
	sppf.Gen(parserDir)
}
