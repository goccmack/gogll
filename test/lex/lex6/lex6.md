Test token reference in lex rule.
```
package "github.com/goccmack/gogll/v3/test/lex/lex6"

id_start : '[\p{L}\p{Nl}\p{Other_ID_Start}-\p{Pattern_Syntax}-\p{Pattern_White_Space}]' ; 

ptrn_stx : '<' ;

```