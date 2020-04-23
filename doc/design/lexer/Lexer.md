# Lexer Generator Design 
The figure below shows the set relations between the different Unicode event classes supported by the lexer:

`A ⊂ B` denotes that a subset relation, `A is a subset of B`, is defined and has values `true` and `false`. This relation is transitive: `A ⊂ B ⊂ C => A ⊂ C`

`A ⋂=∅ B` denotes that a disjoint set relation, `A and B are disjoint`, is defined and has values `true` and `false`.

![Lexer_event_class_relations](fig/set_relations.png)
*Lexer Event Class Relations*