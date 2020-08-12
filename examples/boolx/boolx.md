# Grammar for example: boolx

```
package "boolx"

Expr :   var
     |   Expr Op Expr
     ;

```
The second alternate above, `Expr : Expr Op Expr`, is ambiguous and can produce an ambiguous parse forest.
```

var : letter ;

Op : "&" | "|" ; 

```