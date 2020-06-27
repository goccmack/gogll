package g1

import (
	"fmt"
	"strings"
	"testing"

	"g1/lexer"
	"g1/parser"
)

func Test0(t *testing.T) {
	lex := lexer.New([]rune("a + a + a"))
	if g1, err := parser.New(lex).Parse(); err != nil {
		fmt.Println(err)
		t.Fail()
	} else {
		fmt.Println(strings.Join(g1.([]string), " "))
	}
}
