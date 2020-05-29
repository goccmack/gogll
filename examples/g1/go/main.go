package main

import (
	"g1/lexer"
	"g1/parser"
	"g1/parser/bsr"
	"g1/token"

	"fmt"
	"os"
	"time"
)

type ExprType int

const (
	Expr_And ExprType = iota
	Expr_Or
	Expr_Var
)

type Expr struct {
	Type  ExprType
	Var   *token.Token
	Left  *Expr
	Right *Expr
}

func main() {
	lex := lexer.NewFile(os.Args[1])
	start := time.Now()
	bsrSet, errs := parser.Parse(lex)
	if len(errs) > 0 {
		fail(errs)
	}
	fmt.Printf("%d μs\n", time.Now().Sub(start)/time.Microsecond)
	fmt.Printf("%d BSRs\n", len(bsrSet.GetAll()))
	fmt.Println(getExpr(bsrSet))
	fmt.Printf("%d μs elapsed\n", time.Now().Sub(start)/time.Microsecond)
}

func getExpr(set *bsr.Set) *Expr {
	for _, r := range set.GetRoots() {
		if expr := buildExpr(r); expr != nil {
			return expr
		}
	}
	panic("No valid roots")
}

/*
Exp : Exp Op Exp
    | id
    ;
*/
func buildExpr(b bsr.BSR) *Expr {
	/*** Expr :   var ***/
	if b.Alternate() == 1 {
		return &Expr{
			Type: Expr_Var,
			Var:  b.GetTChildI(0),
		}
	}

	/*** Expr : Expr Op Expr ***/
	op := b.GetNTChildI(1). // Op is symbol 1 of the Expr rule
				GetTChildI(0) // The operator token is symbol 0 for both alternates of the Op rule

	// Build the left subexpression Node. The subtree for it may be ambiguous.
	var left *Expr
	// b.GetNTChildrenI(0) returns all the valid BSRs for symbol 0 of the body of the rule.
	for _, le := range b.GetNTChildrenI(0) {
		// Add subexpression if it is valid and has precedence over this expression
		if e := buildExpr(le); e != nil && e.hasPrecedence(op) {
			left = e
			break
		}
	}
	// No valid subexpressions therefore this whole expression is invalid
	if left == nil {
		return nil
	}
	// Do the same for the right subexpression
	var right *Expr
	for _, re := range b.GetNTChildrenI(2) {
		if e := buildExpr(re); e != nil && e.hasPrecedence(op) {
			right = e
			break
		}
	}
	if right == nil {
		return nil
	}

	// return an expression node
	expr := &Expr{
		Left:  left,
		Right: right,
	}
	switch op.LiteralString() {
	case "&":
		expr.Type = Expr_And
	case "|":
		expr.Type = Expr_Or
	default:
		panic("invalid")
	}
	return expr
}

// & > |, & > &, | > |
func (e *Expr) hasPrecedence(op *token.Token) bool {
	switch e.Type {
	case Expr_And:
		return true
	case Expr_Or:
		return op.LiteralString() == "|"
	case Expr_Var:
		return true
	}
	return false
}

func (e *Expr) String() string {
	switch e.Type {
	case Expr_And:
		return fmt.Sprintf("(%s & %s)", e.Left, e.Right)
	case Expr_Or:
		return fmt.Sprintf("(%s | %s)", e.Left, e.Right)
	case Expr_Var:
		return e.Var.LiteralString()
	}
	panic("invalid")
}

func fail(errs []*parser.Error) {
	fmt.Println("Parse Errors:")
	ln := errs[0].Line
	for _, err := range errs {
		if err.Line == ln {
			fmt.Println("  ", err)
		}
	}
	os.Exit(1)
}
