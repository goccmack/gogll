# Grammar for example: boolx

```
package "github.com/goccmack/gogll/examples/boolx"

Expr :   var
     |   Expr Op Expr
     ;

```
The second alternate above, `Expr : Expr Op Expr`, is ambiguous and can produce an ambiguous parse forest.
```

var : letter ;

Op : "&" | "|" ; 

```