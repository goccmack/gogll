package golang

import (
	"gogll/cfg"
	"gogll/gen/golang/goutil"
	"gogll/gen/golang/parser"
	"path"
)

func Gen() {
	goutil.Gen(path.Join(cfg.BaseDir, "goutil"))
	parser.Gen(path.Join(cfg.BaseDir, "parser"))
}
