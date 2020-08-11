mod lexer;
mod parser;
mod token;

use std::rc::Rc;

fn main() {
    let input1: Rc<Vec<char>> = Rc::new("aname 123".chars().collect());
    let input2: Rc<Vec<char>> = Rc::new("123".chars().collect());

    test1(input1);
    test2(input2);
}

fn test1(input: Rc<Vec<char>>) {
    let lex = lexer::Lexer::new(input);
    let (bsr_set, errs) = parser::parse(lex.clone());
    if errs.len() != 0 {
        panic!()
    }

    // A1 : Name int ;
    let root = bsr_set.get_root();

    // Name : name | empty ;
    let nm = bsr_set.get_nt_child_i(root.clone(), 0);
    if bsr_set.alternate(nm.clone()) != 0 {
        panic!()
    }
    let nmt = bsr_set.get_t_child_i(nm.clone(), 0);
    if nmt.literal_string() != "aname".to_string() {
        panic!();
    }

    let intt = bsr_set.get_t_child_i(root.clone(), 1);
    if intt.literal_string() != "123".to_string() {
        panic!();
    }
}

fn test2(input: Rc<Vec<char>>) {
    let lex = lexer::Lexer::new(input);
    let (bsr_set, errs) = parser::parse(lex.clone());
    if errs.len() != 0 {
        panic!()
    }

    // A1 : Name int ;
    let root = bsr_set.get_root();

    // Name : name | empty ;
    let nm = bsr_set.get_nt_child_i(root.clone(), 0);
    if bsr_set.alternate(nm.clone()) != 1 {
        panic!()
    }

    let intt = bsr_set.get_t_child_i(root.clone(), 1);
    if intt.literal_string() != "123".to_string() {
        panic!(intt.literal_string());
    }
}

