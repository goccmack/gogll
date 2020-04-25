# Example: names

```
package "names"

name : letter { letter | number | '_' } ;
qualifiedName : letter {letter|number|'_'} <'.' <letter|number|'_'>> ;

```