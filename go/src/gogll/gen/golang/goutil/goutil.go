package goutil

import (
	"gogll/gen/golang/goutil/bsr"
	"gogll/gen/golang/goutil/stringset"
	"path/filepath"
)

func Gen(utilDir string) {
	bsr.Gen(filepath.Join(utilDir, "bsr", "bsr.go"))
	stringset.Gen(filepath.Join(utilDir, "stringset", "stringset.go"))
}
