Grammar Γ₁ from:
Derivation representation using binary subtree sets
Elizabeth Scott, Adrian Johnstone, L. Thomas van Binsbergen
Science of Computer Programming 175 (2019

```
package "gamma1" 

S : "a" A B | "a" A "b" ;
A : "a" | "c" | empty ;
B : "b" | B "c" | empty ;
```