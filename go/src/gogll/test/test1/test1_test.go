package test1

import(
	"testing"
	"gogll/test/test1/parser"
	"gogll/test/test1/goutil/bsr"
)

const src =`//C
 1`

func Test1(t *testing.T) {
	err := parser.Parse([]byte(src))
	bsr.Dump()
	if err != nil {
		panic(err)
	}
}