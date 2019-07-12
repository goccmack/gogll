package golang

import (
	"gogll/cfg"
	"gogll/gen/golang/goutil"
	"gogll/gen/golang/parser"
	"path"
)

func Gen() {
	parser.Gen(path.Join(cfg.BaseDir, "parser"))
	goutil.Gen(path.Join(cfg.BaseDir, "goutil"))
}
