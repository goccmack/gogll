# Grammar for logs parser
```
package "logs"

sap : <letter> { '.' <letter> } ':' <number> ;

ip : <number> '.' <number> '.' <number> '.' <number> ;

name : ( '-' | letter { letter | number | '_' } { '.' letter { letter | number | '_' } } ) ;

timestamp : '[' <number> '/' <letter>  '/' <number> ':' <number> ':' <number> ':' <number> ' ' ('+'|'-') <number> ']' ;

string : '"' { not "\"" } '"' ;

number1 : <number> ;

Lines : Line | Lines Line ;

Line : sap ip name name timestamp string number1 number1 string string ;
```
