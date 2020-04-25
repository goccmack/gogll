# Gogll v3

[Copyright 2019 Marius Ackerman](License.txt)

This document contains BNF specification for gogll V3. 

```
package "github.com/goccmack/gogll"

GoGLL : Package Rules ;

Package : "package" string_lit ;

Rules
    :   Rule            
    |   Rule Rules  
    ;

Rule : LexRule | SyntaxRule ;
```

# Lexical Definitions
```
nt : upcase <letter|number|'_'> ;
tokid : lowcase <letter|number|'_'> ; 

char_lit : '\'' (not "\\'" | '\\' any "\\'nrt") '\'' ;
string_lit : '"' {not "\\\"" | '\\' any "\\\"nrt"} '"' ;

```
# Lex Rules
```
LexRule : TokID ":" RegExp ";" ;

RegExp : LexSymbol | LexSymbol RegExp ;

LexSymbol : "." | "any" string_lit | char_lit | LexBracket | "not" string_lit | UnicodeClass ;

LexBracket : LexGroup | LexOptional | LexZeroOrMore | LexOneOrMore ;

LexGroup : "(" LexAlternates ")" ;

LexOptional : "[" LexAlternates "]" ;

LexZeroOrMore : "{" LexAlternates "}" ;

LexOneOrMore : "<" LexAlternates ">" ;

LexAlternates : RegExp | RegExp "|" LexAlternates ;

UnicodeClass : "letter" | "upcase" | "lowcase" | "number" ;

```

# Syntax Rules
```
SyntaxRule : NT ":" SyntaxAlternates ";"  ;

NT : nt  ;

SyntaxAlternates
    :   SyntaxAlternate                   
    |   SyntaxAlternate "|" SyntaxAlternates    
    ;

SyntaxAlternate
    :   SyntaxSymbols                     
    |   "empty"                     
    ;

SyntaxSymbols
    :   SyntaxSymbol                      
    |   SyntaxSymbol SyntaxSymbols              
    ;

SyntaxSymbol : NT | TokID | string_lit ;

TokID : tokid ;

```
# Builtin tokens
-   `.` accepts any character
-   `any String` accepts any character that is an element of `String`
-   `letter` accepts any character from the Unicode letter category
-   `number` accepts any character from the Unicode number category
-   `space` accepts any Unicode white space character
-   `upcase` accepts any upper case letter
-   `lowcase` accepts any lower case letter
-   `not String` accepts any character that is not an element of `String`

