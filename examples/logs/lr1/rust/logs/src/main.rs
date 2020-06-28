mod ast;
mod lexer;
mod parser;
mod token;

const INPUT_FILE: &str = "../../../../../../logs/data/allvhosts.log";

fn main() {
    let lex = lexer::Lexer::new_file(&INPUT_FILE.to_string()).unwrap();
    if let ast::Node::Lines(lns) = parser::Parser::new(lex).parse().unwrap() {
        println!("{} lines",lns.len());
    } else {
        panic!()
    };
}
