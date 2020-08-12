package boolx

import (
	"fmt"
	"os"
	"testing"

	"boolx/lexer"
	"boolx/parser"
	"boolx/parser/bsr"
	"boolx/token"
)

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

const t1Src = `a | b & c | d & e`

func Test1(t *testing.T) {
	bsrSet, errs := parser.Parse(lexer.New([]rune(t1Src)))
	if len(errs) > 0 {
		fail(errs)
	}

	if bsrSet == nil {
		panic("BSRSet == nil")
	}

	for i, r := range bsrSet.GetRoots() {
		fmt.Printf("%d: %s\n", i, buildExpr(r))
	}
}

/*
Expr :   var
     |   Expr Op Expr
     ;
Op : "&" | "|" ;
*/
func buildExpr(b bsr.BSR) *Expr {
	/*** Expr :   var ***/
	if b.Alternate() == 0 {
		return &Expr{
			Type: Expr_Var,
			Var:  b.GetTChildI(0),
		}
	}

	/*** Expr : Expr Op Expr ***/
	op := b.GetNTChildI(1). // Op is symbol 1 of the Expr rule
				GetTChildI(0) // The operator token is symbol 0 for both alternates of the Op rule

	// Build the left subexpression Node. The subtree for it may be ambiguous.
	left := []*Expr{}
	// b.GetNTChildrenI(0) returns all the valid BSRs for symbol 0 of the body of the rule.
	for _, le := range b.GetNTChildrenI(0) {
		// Add subexpression if it is valid and has precedence over this expression
		if e := buildExpr(le); e != nil && hasPrecedence(e, op) {
			left = append(left, e)
		}
	}
	// No valid subexpressions therefore this whole expression is invalid
	if len(left) == 0 {
		return nil
	}
	// Belts and braces
	if len(left) > 1 {
		panic(fmt.Sprintf("%s has %d left children", b, len(left)))
	}
	// Do the same for the right subexpression
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

	// return an expression node
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
