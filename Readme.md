![Apache 2.0 License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)
[![Build Status](https://github.com/goccmack/gogll/workflows/build/badge.svg)](https://github.com/goccmack/gogll/actions)

Copyright 2019 Marius Ackerman. 

# Note
This version does not support Rust. Please use v3.2.0 for Rust or log an issue if you need the features of this version in Rust.

# GoGLL
Gogll generates a GLL or LR(1) parser and FSA-based lexer for any context-free grammar. 
The generated code is Go or Rust.

[Click here](https://goccmack.github.io/posts/2020-05-31_gogll/) for an introduction
to GLL.

See the [LR(1) documentation](doc/lr1/Readme.md) for generating LR(1) parsers.

The generated GLL parser is a clustered nonterminal parser (CNP) following 
[[Scott et al 2019](#Scott-et-al-2019)]. 
CNP is a version of generalised LL parsing (GLL) 
[[Scott & Johnstone 2016](#Scott-et-al-2016)]. 
GLL parsers can parse all context free (CF) languages.

The generated LR(1) parser is a Pager's PGM or Knuth's original LR(1) 
machine [[Pager 1977](Pager-1977)].

The generated lexer is a linear-time finite state automaton FSA 
[[Grune et al 2012](#Grune-et-al-2012)].
The lexer ignores whitespace.

Gogll accepts grammars in markdown files, which is very useful for documenting the grammar.
For example: see
[gogll's own grammar](gogll.md).

GLL has worst-case cubic time and space complexity but linear complexity for all 
LL productions [[Scott et al 2016](#Scott-et-al-2016)]. 
[See here](https://goccmack.github.io/posts/2020-05-31_gogll/) for space and
CPU time measurements of an extreme ambiguous example.
For comparable grammars tested so far gogll produces faster lexers and parsers 
than [gocc](https://github.com/goccmack/gocc) (FSA/LR-1).

# News
## 2022-10-11
SPPF extraction added to the generated code. See [boolx example](examples/boolx/SPPF.md)
## 2022-08-09
Gogll is used to build DAU [DASL](https://dau-technology.github.io/dau-blog/post/2022-08-02-dasl/)

## 2020-08-12
From v3.2.0 gogll supports tokens that can be suppressed by the lexer. This is useful, for example, to implement code comments. See [example](examples/comments/comments.md).

## 2020-06-28
1. Gogll now also generates LR(1) parsers. It supports 
_Pager's Practical General Method, weak compatibility_ as well as 
_Knuths original LR(1) machine_ for
comparison. Pager's PGM generates LR(1) parser tables similar is size to LALR. 
The option to generate a Knuth LR(1) machine is provide for reference.  
See [LR(1) documentation](doc/lr1/Readme.md) for details.

2. Please note that the `-t <target>` option has been replace by `-go` and `-rust`.
See see [usage](#Usage) below.

## 2020-06-01
[See](https://goccmack.github.io/posts/2020-05-31_gogll/) for an introduction
to GLL and a performance comparison of the generated Go and Rust code parsers.

## 2020-05-22
GoGLL v3.1 generates Rust as well as Go parsers with similar performance:

|| Lexer | Parser | Build
|---|---|---|---|
Go | 119 μs | 1324 μs | 0.124s
Rust | 71 μs | 1297 μs | 2.932s

1. The duration was averaged over 1000 repetitions.
1. Build time was measures with the time command.
    1. For Rust: `time cargo build --release`
    2. For Go: `time go build`

See [examples/rust](examples/rust/Readme.md) for the Rust and Go programs used 
for this comparison.

Use gogll's target option to generate a Rust lexer/parser: `-t rust` (see [usage](#Usage) below). 
Gogll generates Go code by default.


## 2020-04-24
1. GoGLL now generates a linear-time FSA lexer matching the CNP parser.
1. This version of *GoGLL is faster than gocc*. It compiles a sample grammar in  
0.074 s, which GoCC compiles in 0.118 s. Gogll compiles itself in 0.041s.

# 

# Benefits and disadvantages of GLL and LR(1)
GLL is a parsing technique that can handle any context-free (CF) language. GLL has
worst case cubic time and space complexity.

LR(1) handles a subset of the context-free languages that can be parsed bottom-up
with one token look-ahead. LR(1) has linear time complexity and its table driven
parser is very efficient. Pager's _Practical General Method_ (PGM) combines
compatible states as they are generated, keeping the state space small.

A GLL parser has more expensive bookkeeping than an LR(1) parser, making the 
LR(1) parser more efficient for parsing very large inputs.

## When to use GLL
1. When the CF grammar that best expresses the problem is not LR(1).
2. When the LR(1) parser has more than a few conflicts that require additional
language symbols or complex grammar refactorisation to resolve.
3. The inputs to be parsed are not too big. GLL works very well for DSLs or 
programming languages.

## When to use LR(1)
1. When the language can be expressed as an LR(1) grammar. A grammar is LR(1) if 
gogll can generate a conflict-free parser for it.
2. When the input is very big, for example: log files containing tens of thousands
of lines.

# Motivation for a separate lexer
The following observations were made while using GoGLLv2 on a couple of projects.

* Most of the ambiguity in grammars were generated by the lexical rules.
* Handling token separation explicitly produces messy, hard to maintain grammars.
* Most of a grammar input file is whitespace, which together with the additional 
ambiguity introduced by the lexical rules, causes most of the parse time in a 
scannerless parser.
* Writing good markdown with the grammar produced slow compilations.

# Input Symbols, Markdown Files
Gogll and lexers generated by gogll accept UTF-8 input strings, which may be in 
a markdown file or a plain text file.

If the input is a markdown file gogll and lexers generated by gogll treat all 
text outside markdown code blocks as whitespace. Markdown code blocks are 
delimited by triple backticks. See [gogll.md](gogll.md) for an example.

# Gogll Grammar
Gogll v3 has a BNF grammar. See [gogll.md](gogll.md)

# Installation
1. Install Go from [https://golang.org](https://golang.org)
1. `go install github.com/goccmack/gogll/v3@latest` or 
1. Clone this repository and run `go install` in the root of the directory where
it is installed.

# Usage
Enter `gogll -h` or `gogll` for the following help:

```
use: gogll -h
    for help, or

use: gogll -version
    to display the version of goggl, or

use: gogll [-a][-v] [-CPUProf] [-o <out dir>] [-go] [-rust] [-gll] [-pager] [-knuth] [-resolve_conflicts] <source file>
    to generate a lexer and parser.

    <source file>: Mandatory. Name of the source file to be processed. 
        If the file extension is ".md" the bnf is extracted from markdown code 
        segments enclosed in triple backticks.
    
    -a: Optional. Regenerate all files.
        WARNING: This may destroy user editing in the LR(1) AST.
        Default: false
         
    -v: Optional. Produce verbose output, including first and follow sets,
        LR(1) sets and lexer FSA sets.
    
    -o <out dir>: Optional. The directory to which code will be generated.
                  Default: the same directory as <source file>.
                  
    -go: Optional. Generate Go code.
          Default: true, but false if -rust is selected

    -rust: Optional. Generate Rust code.
           Default: false
           
    -gll: Optional. Generate a GLL parser.
          Default true. False if -knuth or -pager is selected.
                  
    -knuth: Optional. Generate a Knuth LR(1) parser
            Default false

    -pager: Optional. Generate a Pager PGM LR(1) parser.
            Default false

    -resolve_conflicts: Optional. Automatically resolve LR(1) conflicts.
            Default: false. Only used when generating LR(1) parsers.
    
    -bs: Optional. Print BSR statistics (GLL only).
    
    -CPUProf : Optional. Generate a CPU profile. Default false.
        The generated CPU profile is in <cpu.prof>. 
        Use "go tool pprof cpu.prof" to analyse the profile.
```

# Using the generated lexer and parser
1. Create a lexer:  
From an `[]rune`:
```
	lexer.New(input []rune) *Lexer
```
  or from a file. If the file extension us `.md` the lexer will 
  treat all text outside the markdown code blocks as whitespace.
```
	lexer.NewFile(fname string) *Lexer
```
2. Parse the lexer:  
```
	if err, errs := parser.Parse(lex); err != nil {...}
```
3. Check for ambiguities in the parse forest
```
	if bsr.IsAmbiguous() {
		fmt.Println("Error: Ambiguous parse forest")
		bsr.ReportAmbiguous()
		os.Exit(1)
	}
```
Ambiguous BSRs must be resolved by walking the parse forest and ignoring
unwanted children of ambiguous NTs (see [Complete Example](#Complete-Example)).
4. Use the disambiguated parse tree for the further stages of compilation. 
For example, see gogll's [AST builder](ast/build.go).

<a name="Complete-Example"></a>
# Complete Example
The code of following example can be found at [examples/boolx](examples/boolx/boolx.md). 
The example has the following grammar: [boolx.md](examples/boolx/boolx.md), which generates boolean expressions such as: `a | b & c | d & e`:

```
package "github.com/goccmack/gogll/examples/boolx"

Expr :   var
     |   Expr Op Expr
     ;

var : letter ;

Op : "&" | "|" ; 
```
The second alternate above, `Expr : Expr Op Expr`, is ambiguous and can produce an ambiguous parse forest.
The grammar does not enforce operator precedence, 
this has to be done during semantic analysis.

The grammar is compiled by the following command:
```
gogll examples/boolx/boolx.md
```

The test file, [boolx_test.go](examples/boolx/boolx_test.go) shows the steps
required to parse an input string and produce a disambiguated abstract syntax tree:

```
const t1Src = `a | b & c | d & e`

func Test1(t *testing.T) {
```
1. Create a lexer from the input string and parse. Fail if there are parse errors.
```
	if err, errs := parser.Parse(lexer.New([]rune(t1Src))); err != nil {
		fail(errs)
	}

```
2. Build an abstract syntax tree for each root of the parse forest and print them.
```
	for i, r := range bsr.GetRoots() {
		fmt.Printf("%d: %s\n", i, buildExpr(r))
	}
}
```
The input string produces an ambiguous parse forest, which is partially 
disambiguated by applying operator precedence.
We get the following output from this test:
```
> go test -v ./examples/boolx
=== RUN   Test1
0: (a | ((b & c) | (d & e)))
1: <nil>
2: <nil>
3: ((a | (b & c)) | (d & e))
--- PASS: Test1 (0.00s)
PASS
```
The output shows that the parse forest has 4 roots, 2 of which produce valid ASTs 
after disambiguation. The removed trees are syntactically valid by semantically
invalid because they give `|` precedence over `&`. 
Both the remaining ASTs are syntactically and semantically
valid. The AST encodes operator precedence as shown by the parentheses. 
The choice of which valid AST to use for further processing is application specific.

In this example disambiguation by operator precedence is applied during the
AST build. 

Our AST has only one type of node: `Expr`.
```
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

```
A node can represent a variable (`Type` = `Expr_Var`) or an expression (`Type` = `Expr_Expr`).
If the node represents a variable the field `Var` contains the variable token. 
Otherwise `Op` contains the operator token and `Left` and `Right` contain the nodes
of the sub-expressions.

The AST is constructed recursively from each BSR root by the function, `buildExpr`
in [boolx_test.go](examples/boolx/boolx_test.go).

```
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
```

# Status
* `gogll v3` generates a matching lexer and parser. It generates GLL and LR(1) 
parsers. v3 compiles itself.
v3 is used in a real-world project.
* `gogll v2` had the last vestiges of the bootstrap compiler grammar removed from
its input grammar. v2 compiled itself.
* `gogll v1` was a GLL scannerless parser, which compiled scannerless GLL parsers.
v1 compiled itself.
* `gogll v0` was a bootstrap compiler implemented by a [gocc](https://github.com/goccmack/gocc) lexer and parser.

# Features considered for future implementation
1. Tokens suppressed by the lexer, e.g.: code comments.
1. Better error reporting.
1. Better documentation, including how to traverse the binary subtree representation (BSR [Scott et al 2019](#Scott-et-al-2019)) of the parse forest as well as on disambiguating 
parse forests.
1. Letting the parser direct which tokens to scan [Scott & Johnstone 2019](#Scott-et-al-2019a)

# Documentation
At the moment this document and the [gogll grammar](gogll.md) are the only documentation. Have a look at 
`gogll/examples/ambiguous` for a simple example and also for simple disambiguation.

Alternatively look at `gogll.md` which is the input grammar and also the grammar
from which the `parser` for this version of `gogll` was generated. `gogll/da` disambiguates the parse forest for an input string.

## LR(1)
See the [LR(1) documentation](doc/lr1/Readme.md).

# Changelog
[see](ChangeLog.md)

# Bibliography
<a name=Pager-1977></a>
* [Pager 1977] David Pager   
A Practical General Method for Constructing LR(k) Parsers   
Acta Informatica 7, 1977

<a name=Scott-et-al-2019a></a>
* [Scott & Johnstone 2019] Elizabeth Scott and Adrian Johnstone  
Multiple lexicalisation (a Java based study)  
In: [Proceedings of Software Language Engineering 2019. ACM, 2019. p. 71-82](https://pure.royalholloway.ac.uk/portal/files/34483813/lcnpSubmitFromEASForPure.pdf)

<a name="Scott-et-al-2019"></a>
* [Scott et al 2019] Elizabeth Scott, Adrian Johnstone and L. Thomas van Binsbergen.  
Derivation representation using binary subtree sets.  
In: Science of Computer Programming (175) 2019

<a name="Scott-et-al-2018"></a>
* [Scott & Johnstone 2018] Elizabeth Scott and Adrian Johnstone.   
GLL Syntax Analysers For EBNF Grammars.   
In: [Science of Computer Programming
Volume 166, 15 November 2018](https://pure.royalholloway.ac.uk/portal/en/publications/gll-syntax-analysers-for-ebnf-grammars(58d1ec5e-28df-486a-879e-36d58a9f8abf).html)

<a name="Scott-et-al-2016"></a>
* [Scott & Johnstone 2016] Elizabeth Scott and Adrian Johnstone.   
Structuring the GLL parsing algorithm for performance.   
In: [Science of Computer Programming
Volume 125, 1 September 2016](https://pure.royalholloway.ac.uk/portal/en/publications/structuring-the-gll-parsing-algorithm-for-performance(a95fc020-9918-4f17-a87a-845e2aee12b8).html)

<a name="Afroozeh-et-al-2013"></a>
* [Afroozeh et al 2013] Ali Afroozeh, Mark van den Brand, Adrian Johnstone, Elizabeth Scott, Jurgen Vinju.   
Safe Specification of Operator Precedence Rules.   
In: [Erwig M., Paige R.F., Van Wyk E. (eds) Software Language Engineering. SLE 2013. Lecture Notes in Computer Science, vol 8225. Springer, Cham](https://pure.royalholloway.ac.uk/portal/en/publications/safe-specification-of-operator-precedence-rules(0287d70e-92b8-4204-aafb-15a81de84968).html)

<a name="Grune-et-al-2012"></a>
* [Grune et al 2012] Dick Grune, Kees van Reeuwijk, Henri E. Bal, Ceriel J.H. Jacobs and Koen Langendoen.
Modern Compiler Design. Second Edition.
Springer 2012

<a name="Basten-2012"></a>
* [Basten & Vinju 2012] Basten H.J.S., Vinju J.J. (2012) Parse Forest Diagnostics with Dr. Ambiguity. In: Sloane A., Aßmann U. (eds) Software Language Engineering. SLE 2011. [Lecture Notes in Computer Science, vol 6940. Springer, Berlin, Heidelberg](https://homepages.cwi.nl/~jurgenv/papers/SLE2011-2.pdf)

