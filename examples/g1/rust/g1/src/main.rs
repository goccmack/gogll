mod lexer;
mod parser;
mod token;

use lexer::Lexer;
use parser::bsr;
use token::Token;

use std::rc::Rc;

#[derive(Debug)]
enum Expr {
    And{left: Box<Expr>, right: Box<Expr>},
    Or{left: Box<Expr>, right: Box<Expr>},
    Id(String),
}

#[derive(Clone, Copy, Debug, PartialEq, Eq)]
enum Op {
    And,
    Or,
}

fn main() {
    let input_file = &std::env::args().collect::<Vec<String>>()[1];
    let lex = Lexer::new_file(&input_file).unwrap();
    let start = std::time::SystemTime::now();
    let (bsr_set, errs) = parser::parse(lex.clone());
    if errs.len() > 0 {
        fail(&errs)
    }
    println!("{} μs parse", start.elapsed().unwrap().as_micros());
    println!("{} BSRs", bsr_set.get_all().len());
    println!("{}", get_exp(&bsr_set).to_string());
    println!("{} μs elapsed", start.elapsed().unwrap().as_micros());
}

fn get_exp(set: &parser::bsr::Set) -> Box<Expr> {
    for r in set.get_roots() {
        if let Some(e) = exp(set, r.clone()) {
            return e;
        }
    }
    panic!("No valid roots found");
}

/*
Exp : Exp Op Exp
    | id
    ;
*/
fn exp(set: &bsr::Set, b: Rc<bsr::BSR>) -> Option<Box<Expr>> {
    if set.alternate(b.clone()) == 1 {
        return Some(id(set, b));
    }

    let op: Op = op(set, set.get_nt_child_i(b.clone(), 1));

    let mut left: Option<Box<Expr>> = None;
    for b1 in set.get_nt_children_i(b.clone(), 0) {
        if let Some(e) = exp(set, b1.clone()) {
            if e.has_precedence(op) {
                left = Some(e);
                break;
            }
        }
    }
    if left.is_none() {
        return None;
    }

    let mut right: Option<Box<Expr>> = None;
    for b1 in set.get_nt_children_i(b.clone(), 2) {
        if let Some(e) = exp(set, b1.clone()) {
            if e.has_precedence(op) {
                right = Some(e);
                break;
            }
        }
    }
    if right.is_none() {
        return None;
    }
    // If there are multiple valid left or right sub-expressions we pick
    // the last one
    let l: Box<Expr> = left.unwrap();
    let r: Box<Expr> = right.unwrap();
    match op {
        Op::And => Some(Box::new(Expr::And{left:l, right: r})),
        Op::Or => Some(Box::new(Expr::Or{left:l, right: r})),
    }
}

// Exp : id
fn id(set: &bsr::Set, b: Rc<bsr::BSR>) -> Box<Expr> {
    let id: Rc<Token> = set.get_t_child_i(b, 0);
    let id_str: String = id.literal().iter().collect();
    Box::new(Expr::Id(id_str))
}

// Op : "&" | "|" ;
fn op(set: &bsr::Set, b: Rc<bsr::BSR>) -> Op {
    match set.alternate(b.clone()) {
        0 => Op::And,
        1 => Op::Or,
        _ => panic!("Invalid alternate {}", set.alternate(b)),
    }
}

// Print all the errors with the same line number as errs[0] and exit(1)
fn fail(errs: &Vec<Box<parser::Error>>) {
    println!("Parse Errors:");
    let ln = errs[0].line;
    for err in errs {
        if err.line == ln {
            println!("  {}", err);
        }
    }
    std::process::exit(1);
}

impl Expr {
    #[allow(dead_code)]
    fn to_string(&self) -> String {
        match self {
            Expr::And{left, right} => {
                let mut result = "(".to_string();
                result.push_str(&left.to_string());
                result.push_str(&" & ");
                result.push_str(&right.to_string());
                result.push(')');
                return result
            },
            Expr::Or{left, right} => {
                let mut result = "(".to_string();
                result.push_str(&left.to_string());
                result.push_str(&" | ");
                result.push_str(&right.to_string());
                result.push(')');
                return result
            },
            Expr::Id(id) => id.clone(),
        }
    }

    fn has_precedence(&self, op: Op) -> bool {
        match self {
            Expr::And{left:_,right:_} => true, 
            Expr::Or{left:_,right:_} => return op == Op::Or,
            Expr::Id(_) => true,
        }
    }
} // impl Expr

