package parser

import(
	"testing"
	"fmt"
	"gogll/test/AAc/parser/sppf"
)

const src = "c"

func Test(t *testing.T) {
	Parse([]byte(src))
	fmt.Println(sppf.Dot())
}