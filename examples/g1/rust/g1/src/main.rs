mod lexer;
mod parser;
mod token;

use lexer::Lexer;
use std::rc::Rc;

fn main() {
    let input: Rc<Vec<char>> = Rc::new("a | b & c | d".chars().collect());
    let lex = Lexer::new(input);
    let (_bsr_set, errs) = parser::parse(lex.clone());
    if errs.len() > 0 {
        fail(&errs)
    }
}

fn fail(errs: &Vec<Box<parser::Error>>) {
    println!("Parse Errors:");
    let ln = errs[0].line;
    for err in errs {
        if err.line == ln {
            println!("  {}", err);
        }
    }
    std::process::exit(-1);
}
