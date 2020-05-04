# Grammar for example: bool

```
package "github.com/goccmack/gogll/examples/boolx"

Expr :   var
     |   Expr Op Expr
     ;

var : letter ;

Op : "&" | "|" ; 

```