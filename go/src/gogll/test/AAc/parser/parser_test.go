package parser

import(
	"testing"
	"gogll/test/AAc/gen/dot"
)

const src = "c"

func Test(t *testing.T) {
	Parse([]byte(src))
	dot.Gen(".")
}