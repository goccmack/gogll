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

    // Disambiguate parse forest and print resulting expression
    println!("{}", get_expr(&bsr_set).to_string());

    println!("{} μs elapsed", start.elapsed().unwrap().as_micros());
}

// Return the first logically valid expression
fn get_expr(set: &parser::bsr::Set) -> Box<Expr> {
    for r in set.get_roots() {
        if let Some(e) = build_expr(set, r.clone()) {
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
fn build_expr(set: &bsr::Set, b: Rc<bsr::BSR>) -> Option<Box<Expr>> {
    // Alternate 1 of the rule
    if set.alternate(b.clone()) == 1 {
        return Some(build_id(set, b));
    }

    // Get Op: symbol 1 in the body of alternate 0
    let op: Op = build_op(set, set.get_nt_child_i(b.clone(), 1));

    // Build the left subexpression Node. The subtree for it may be ambiguous.
    // Pick the first subexpression that has operator precedence over this expression.
    let mut left: Option<Box<Expr>> = None;
    // set.get_nt_children_i(b, 0) returns all the valid subtrees for symbol 0 of
    // Exp : Exp Op Exp
    for b1 in set.get_nt_children_i(b.clone(), 0) {
		// Pick the first subexpression with operator precedence over this expression
        if let Some(e) = build_expr(set, b1.clone()) {
            if e.has_precedence(op) {
                left = Some(e);
                break;
            }
        }
    }
	// No valid subexpressions therefore this whole expression is invalid
    if left.is_none() {
        return None;
    }

	// Do the same for the right subexpression
    let mut right: Option<Box<Expr>> = None;
    for b1 in set.get_nt_children_i(b.clone(), 2) {
        if let Some(e) = build_expr(set, b1.clone()) {
            if e.has_precedence(op) {
                right = Some(e);
                break;
            }
        }
    }
    if right.is_none() {
        return None;
    }

	// return an expression node
    let l: Box<Expr> = left.unwrap();
    let r: Box<Expr> = right.unwrap();

    match op {
        Op::And => Some(Box::new(Expr::And{left:l, right: r})),
        Op::Or => Some(Box::new(Expr::Or{left:l, right: r})),
    }
}

// Exp : id
fn build_id(set: &bsr::Set, b: Rc<bsr::BSR>) -> Box<Expr> {
    let id: Rc<Token> = set.get_t_child_i(b.clone(), 0);
    let id_str: String = id.literal().iter().collect();
    Box::new(Expr::Id(id_str))
}

// Op : "&" | "|" ;
fn build_op(set: &bsr::Set, b: Rc<bsr::BSR>) -> Op {
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

    // id > & > |
    fn has_precedence(&self, op: Op) -> bool {
        match self {
            Expr::And{left:_,right:_} => true, 
            Expr::Or{left:_,right:_} => return op == Op::Or,
            Expr::Id(_) => true,
        }
    }
} // impl Expr

