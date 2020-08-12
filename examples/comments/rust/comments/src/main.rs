mod lexer;
mod token;

fn main() {
    let toks = &lexer::Lexer::new_file(&"../../test.txt".to_string()).unwrap().tokens;
    for t in toks.iter() {
        println!("{}", t.clone());
    }
    if toks[2].typ != token::Type::EOF {
        panic!()
    }
    let tok = toks[0].clone();
    if tok.id() != "name" || tok.literal_string() != "name1" {
        println!("id {}, lit {}", tok.id(), tok.literal_string());
        panic!(tok.id().to_string());
    }
    let tok = toks[1].clone();
    if tok.id() != "name" || tok.literal_string() != "name2" {
        panic!(tok.id().to_string());
    }
    println!("OK")
}
