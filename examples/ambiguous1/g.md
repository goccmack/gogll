# Example: Ambiguous grammar

```
package "gogll/examples/ambiguous1"

*S : A S | B S | emptyAlt ;

A : letter ;

B : letter ;
```