# Gogll v3.0

[Copyright 2019 Marius Ackerman](License.txt)

This document contains BNF specification for gogll V3. 

# Lexical Definitions
```
package "github.com/goccmack/gogll"

GoGLL : Package Rules ;

Package : "package" string_lit ;

Rules
    :   Rule            
    |   Rule Rules
    ;

Rule : NT ":" Alternates ";"  ;

NT : id  ;

Alternates
    :   Alternate                   
    |   Alternate "|" Alternates    
    ;

Alternate
    :   Symbols                     
    |   "empty"                     
    ;

Symbols
    :   Symbol                      
    |   Symbol Symbols              
    ;

Symbol : NT | TokID | string_lit ;

TokID : id ;

```

-   `any` accepts any character
-   `anyof String` accepts any character that is an element of `String`
-   `letter` accepts any character from the Unicode letter category
-   `number` accepts any character from the Unicode number category
-   `space` accepts any Unicode white space character
-   `upcase` accepts any upper case letter
-   `lowcase` accepts any lower case letter
-   `not "String"` accepts any character that is not an element of `String`

