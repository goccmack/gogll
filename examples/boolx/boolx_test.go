package boolx

import (
	"fmt"
	"os"
	"testing"

	"github.com/goccmack/gogll/examples/boolx/lexer"
	"github.com/goccmack/gogll/examples/boolx/parser"
	"github.com/goccmack/gogll/examples/boolx/parser/bsr"
	"github.com/goccmack/gogll/examples/boolx/token"
)

const t1Src = `a | b & c | d & e`

type ExprType int

const (
	Expr_Var ExprType = iota
	Expr_Expr
)

type Expr struct {
	Type  ExprType
	Var   *token.Token
	Op    *token.Token
	Left  *Expr
	Right *Expr
}

func Test1(t *testing.T) {
	if err, errs := parser.Parse(lexer.New([]rune(t1Src))); err != nil {
		fail(errs)
	}

	for _, r := range bsr.GetRoots() {
		if e := buildExpr(r); e != nil {
			fmt.Println(e)
		}
	}
}

/*
Expr :   var
     |   Expr Op Expr
     ;
Op : "&" | "|" ;
*/
func buildExpr(b bsr.BSR) *Expr {
	if b.Alternate() == 0 {
		return &Expr{
			Type: Expr_Var,
			Var:  b.GetTChildI(0),
		}
	}
	op := b.GetNTChildI(1).GetTChildI(0)
	left := []*Expr{}
	for _, le := range b.GetNTChildrenI(0) {
		if e := buildExpr(le); e != nil && hasPrecedence(e, op) {
			left = append(left, e)
		}
	}
	if len(left) == 0 {
		return nil
	}
	if len(left) > 1 {
		panic(fmt.Sprintf("%s has %d left children", b, len(left)))
	}
	right := []*Expr{}
	for _, le := range b.GetNTChildrenI(2) {
		if e := buildExpr(le); e != nil && hasPrecedence(e, op) {
			right = append(right, e)
		}
	}
	if len(right) == 0 {
		return nil
	}
	if len(right) > 1 {
		panic(fmt.Sprintf("%s has %d right children", b, len(right)))
	}

	return &Expr{
		Type:  Expr_Expr,
		Op:    op,
		Left:  left[0],
		Right: right[0],
	}
}

// & > |, & > &, | > |
func hasPrecedence(e *Expr, op *token.Token) bool {
	if e.Type == Expr_Var {
		return true
	}
	return e.Op.LiteralString() == "&" || op.LiteralString() == "|"
}

func (e *Expr) String() string {
	if e.Type == Expr_Var {
		return e.Var.LiteralString()
	}
	return fmt.Sprintf("(%s %s %s)", e.Left, e.Op.LiteralString(), e.Right)
}

func fail(errs []*parser.Error) {
	ln := errs[0].Line
	fmt.Println("Parse Error:")
	for _, e := range errs {
		if e.Line == ln {
			fmt.Println(e)
		}
	}
	os.Exit(1)
}
