mod ast;
mod lexer;
mod parser;
mod token;

use std::time::SystemTime;

const INPUT_FILE: &str = "../../../data/test.log";

fn main() {
    let start = SystemTime::now();
    let lex = lexer::Lexer::new_file(&INPUT_FILE.to_string()).unwrap();
    let lexDone = SystemTime::now();
    if let ast::Node::Lines(lns) = parser::Parser::new(lex).parse().unwrap() {
    } else {
        panic!()
    };
    let parseDone = SystemTime::now();
    println!("Lexer duration {} ms", lexDone.duration_since(start).unwrap().as_millis());
    println!("Parse duration {} ms", parseDone.duration_since(lexDone).unwrap().as_millis());
}


// Test1
// fn main() {
//     let REP = 1000*1000*1000;
//     let start = SystemTime::now();
//     for _i in 1..REP {
//         let _ = lexer::NEXT_STATE[0]('[');
//     }
//     println!("Test1 duration {} ms", start.elapsed().unwrap().as_millis());
// }

// Test2
// fn main() {
//     let REP = 100*1000;
//     let mut v: Vec<usize> = Vec::with_capacity(2048);
//     let start = SystemTime::now();
//     for i in 0..REP {
//         v.push(i);
//     }
//     println!("Test2 duration {} mus", start.elapsed().unwrap().as_micros();
//     println!("len(v) {}", v.len());
// }

// Test3
// fn main() {
//     let REP = 100*1000;
//     let start = SystemTime::now();
//     for i in 0..REP {
//         '9'.is_alphabetic();
//     }
//     println!("Test3 duration {} mus", start.elapsed().unwrap().as_micros());
// }

// Test4
// fn main() {
//     let REP = 100*1000*1000;
//     let start = SystemTime::now();
//     for _i in 1..REP {
//         let _ = lexer::NEXT_STATE[1]('[');
//     }
//     println!("Test4 duration {} ms", start.elapsed().unwrap().as_millis());
// }

// Test5
// fn main() {
//     let REP = 1000*1000*1000;
//     let ptrn: Vec<char> = vec!['"'];
//     let start = SystemTime::now();
//     let mut res: Vec<bool> = Vec::with_capacity(REP);
//     for _ in 0..REP {
//         for i in 0..ptrn.len() {
//             res.push('a' == ptrn[i]);
//         }
//     }
//     println!("Test5 duration {} ms", start.elapsed().unwrap().as_millis());
//     println!("{}", res[0]);
// }
