# Test back tick
```
package "lex4"

Exp : Exp Op Exp
    | id
    ;

Op : "&" | "|" ;

id : letter <letter | number | '`'> ;
```