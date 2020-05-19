# Rust example and benchmark against Go
This is an example demonstrating how to generate a Rust lexer and parser for 
the [gogll grammar](gogll.md), as well as a demonstration of walking the BSR set in Rust.

# Generating the lexer and parser
The following rule from [examples/rust/makefile](makefile) shows
how to use gogll generate the Go and Rust code.
```
gogll: 
	gogll -t rust gogll.md
	gogll -o go gogll.md
```
# Running
```
cd examples/rust
make
```

# Rust example
The main module declares the lexer and parser modules: `token`, `lexer` and `parser`.
```
mod token;
mod lexer;
mod parser;

use lexer::{Lexer};
use parser::bsr;
use parser::bsr::{BSR};

use std::rc::Rc;
use std::time::SystemTime;

// The number of times to repeat lexing and parsing
const N: u128 = 10000;

fn main() {
    let src_file = "gogll.md".to_string();
```

Run the lexer N times for benchmarking
```
    // Repeat lexing N times
    let start_time = SystemTime::now();
    let mut lex = 
        Lexer::new_file(&src_file).expect(&format!("Error lexing {}", src_file));
    for _ in 1..N-1 {
        lex = 
            Lexer::new_file(&src_file).expect(&format!("Error lexing {}", src_file));
    }
    let lex_done = SystemTime::now();
```

Run the parser N times for benchmarking
```
    // Repeat parsing N times
    let (bsr_set, errs) = parser::parse(lex.clone());
    if errs.len() > 0 {
        panic!("Errors");
    }
    for _ in 1..N-1 {
        parser::parse(lex.clone());
        if errs.len() > 0 {
            panic!("Errors");
        }
    }
    let parse_done = SystemTime::now();
```

Check if the parse was ambiguous. It isn't so no disambiguation is required.
```
    // Check BSR set for ambiguity
    if bsr_set.is_ambiguous() {
        panic!("Ambiguous BSR Set")
    }
```
Partially walk the BSR set return by the parser, descending from the root. We 
use a separate function for each grammar rule.
```
    // Partially walk the BSR
    walk_bsr(&bsr_set);
```

Report the average time taken by the lexer and parser.
```
    println!("Lexer took {} μs", 
        lex_done.duration_since(start_time).expect("").as_micros()/N);
    println!("Parser took {} μs", 
        parse_done.duration_since(lex_done).expect("").as_micros()/N);
}
```
## The BSR set walking routines
```
// Walk BSR set from root. Each BSR represents a parsed grammar rule
fn walk_bsr(set: &bsr::Set) {
    walk_gogll(set, set.get_root());
}

// GoGLL : Package Rules ;
fn walk_gogll(set: &bsr::Set, b: Rc<BSR>) {
    walk_package(set, set.get_nt_child_i(b.clone(), 0));
    walk_rules(set, set.get_nt_child_i(b.clone(), 1));
}

// Package : "package" string_lit ;
fn walk_package(set: &bsr::Set, b: Rc<BSR>) {
    let _ = set.get_t_child_i(b.clone(), 1);
}

/*
Rules
    :   Rule            
    |   Rule Rules  
    ;
*/
fn walk_rules(set: &bsr::Set, b: Rc<BSR>) {
    let mut b1 = b.clone();
    while set.alternate(b1.clone()) == 1 {
        walk_rule(set, set.get_nt_child_i(b1.clone(), 0));
        b1 = set.get_nt_child_i(b1.clone(), 1);
    }
    walk_rule(set, set.get_nt_child_i(b1.clone(), 0));
}

// Rule : LexRule | SyntaxRule ;
fn walk_rule(set: &bsr::Set, b: Rc<BSR>) {
    match set.alternate(b.clone()) {
        0 => walk_lex_rule(set, set.get_nt_child_i(b, 0)),
        1 => walk_syntax_rule(set, set.get_nt_child_i(b, 0)),
        _ => panic!()
    }
}

fn walk_lex_rule(_set: &bsr::Set, _b: Rc<BSR>) {
    // do nothing
}

fn walk_syntax_rule(_set: &bsr::Set, _b: Rc<BSR>) {
    // do nothing
}
```
