mod ast; 
mod lexer;
mod parser;
mod token;

use std::rc::Rc;

fn main() {
    let input: Rc<Vec<char>> = Rc::new("a + a + a".chars().collect());
    let lex = lexer::Lexer::new(input.clone());
    match parser::Parser::new(lex).parse() {
        Ok(_) => {
            println!("Ok");
        }
        Err(msg) => {
            println!("Error: {}", msg);
            std::process::exit(1);
        }
    }
}
