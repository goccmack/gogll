# Test empty productions
```
package "empty"

name : letter {letter} ;

int : number {number} ;

A1 : Name int ;

Name : name | empty ;
```