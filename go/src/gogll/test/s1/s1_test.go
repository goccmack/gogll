package s1

import (
	"gogll/test/s1/goutil/bsr"
	"testing"
	"gogll/test/s1/parser"
)

func Test1(t *testing.T) {
	parser.Parse([]byte("d"))
	bsr.Dump()
}