Test token reference in lex rule.
```
package "github.com/goccmack/gogll/v3/test/lex/lex6"

id : lowcase <letter|number|'_'> ; 

ids: <id '-'> ;
```