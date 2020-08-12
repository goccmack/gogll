# Example of code comments

This example shows how the lexer can be instructed to suppress tokens that are 
unwanted in the syntax, for example: code comments.

This grammar specifies a list of names that can be written with any number of c-style line and 
block comments anywhere between the names.

The `!` in front of `!line_comment` and `!block_comment` instructs the lexer to 
suppress those tokens. See the [grammar for details.](../../gogll.md)

```
package "comments"

name : letter {letter | number} ;

!line_comment : '/' '/' {not "\n"} '\n' ;
```
`!line_comment` is a c-style line comment. Everything from the first slash to the end of line
is a comment.
```
!block_comment : '/''*' {not "*"} '*''/' ;
```
`!block_comment` is a c-style block comment. Everything between and including 
`/*` and `*/` is a comment. 