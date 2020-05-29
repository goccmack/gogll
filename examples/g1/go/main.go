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
	Expr_Id
)

type Expr struct {
	Type  ExprType
	Id    *token.Token
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

	// Disambiguate parse forest and print resulting expression
	fmt.Println(getExpr(bsrSet))

	fmt.Printf("%d μs elapsed\n", time.Now().Sub(start)/time.Microsecond)
}

// Return the first logically valid expression
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
	// Alternate 1 of the rule
	if b.Alternate() == 1 {
		return buildID(b)
	}

	// Get Op: symbol 1 in the body of alternate 0
	op := buildOp(b.GetNTChildI(1))

	// Build the left subexpression Node. The subtree for it may be ambiguous.
	// Pick the first subexpression that has operator precedence over this expression.
	var left *Expr
	// b.GetNTChildrenI(0) returns all the valid subtrees for symbol 0 of
	// Exp : Exp Op Exp
	for _, le := range b.GetNTChildrenI(0) {
		// Pick the first subexpression with operator precedence over this expression
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

// Exp : id ;
func buildID(b bsr.BSR) *Expr {
	return &Expr{
		Type: Expr_Id,
		Id:   b.GetTChildI(0),
	}
}

// Op : "&" | "|" ;
func buildOp(b bsr.BSR) *token.Token {
	return b.GetTChildI(0)
}

// id > & > |
func (e *Expr) hasPrecedence(op *token.Token) bool {
	switch e.Type {
	case Expr_And:
		return true
	case Expr_Or:
		return op.LiteralString() == "|"
	case Expr_Id:
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
	case Expr_Id:
		return e.Id.LiteralString()
	}
	panic("invalid")
}

// Print all the errors with the same line number as errs[0] and exit(1)
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
