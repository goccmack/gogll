package comments

import (
	"comments/lexer"
	"comments/token"
	"fmt"
	"testing"
)

func Test1(t *testing.T) {
	toks := lexer.NewFile("../test.txt").Tokens
	fmt.Println(toks)
	if toks[2].Type() != token.EOF {
		t.Fail()
	}
	if toks[0].TypeID() != "name" || toks[0].LiteralString() != "name1" {
		t.Fail()
	}
	if toks[0].TypeID() != "name" || toks[1].LiteralString() != "name2" {
		t.Fail()
	}
}
