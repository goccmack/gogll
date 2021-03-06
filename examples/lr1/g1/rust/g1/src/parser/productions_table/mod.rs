//! Generated by GoGLL. Do not edit.

use crate::ast;
use lazy_static::lazy_static;

pub struct ProdTabEntry {
	pub string: String,
	id: String,
	pub nt_type: usize,
	index: usize,
	pub num_symbols: usize,
	pub reduce_func: fn(Vec<ast::Node>) -> Result<ast::Node, String>,
}

lazy_static! {
    pub static ref PROD_TABLE: Vec<ProdTabEntry> = {
		let mut m: Vec<ProdTabEntry> = Vec::with_capacity(0);
		
        m.push(ProdTabEntry{
			string: "G0 : E1 ;".to_string(),
			id: "G0".to_string(),
			nt_type: 0,
			index: 0,
			num_symbols: 1,
			reduce_func: ast::g_0_0,
				// |x: Vec<ast::Node>| -> Result<ast::Node, String> {
				// 	ast::g_0_0(x[0])
				// },
		});
		
        m.push(ProdTabEntry{
			string: "E1 : E1 + T1 ;".to_string(),
			id: "E1".to_string(),
			nt_type: 0,
			index: 1,
			num_symbols: 3,
			reduce_func: ast::e_1_0,
				// |x: Vec<ast::Node>| -> Result<ast::Node, String> {
				// 	ast::e_1_0(x[0],x[1],x[2])
				// },
		});
		
        m.push(ProdTabEntry{
			string: "E1 : T1 ;".to_string(),
			id: "E1".to_string(),
			nt_type: 0,
			index: 2,
			num_symbols: 1,
			reduce_func: ast::e_1_1,
				// |x: Vec<ast::Node>| -> Result<ast::Node, String> {
				// 	ast::e_1_1(x[0])
				// },
		});
		
        m.push(ProdTabEntry{
			string: "T1 : a ;".to_string(),
			id: "T1".to_string(),
			nt_type: 1,
			index: 3,
			num_symbols: 1,
			reduce_func: ast::t_1_0,
				// |x: Vec<ast::Node>| -> Result<ast::Node, String> {
				// 	ast::t_1_0(x[0])
				// },
		});
		
        m
    };
}
