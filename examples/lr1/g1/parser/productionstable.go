
package parser

import(
    "g1/ast"
)

type (
	//TODO: change type and variable names to be consistent with other tables
	ProdTab      [numProductions]ProdTabEntry
	ProdTabEntry struct {
		String     string
		Id         string
		NTType     int
		Index int
		NumSymbols int
		ReduceFunc func([]interface{}) (interface{}, error)
	}
)

var productionsTable = ProdTab {
	ProdTabEntry{
		String: `G0 : E1 ;`,
		Id: "G0",
		NTType: 0,
		Index: 0,
		NumSymbols: 1,
		ReduceFunc: func(X []interface{}) (interface{}, error) {
            return ast.G00(X[0])
		},
	},
	ProdTabEntry{
		String: `E1 : E1 + T1 ;`,
		Id: "E1",
		NTType: 0,
		Index: 1,
		NumSymbols: 3,
		ReduceFunc: func(X []interface{}) (interface{}, error) {
            return ast.E10(X[0],X[1],X[2])
		},
	},
	ProdTabEntry{
		String: `E1 : T1 ;`,
		Id: "E1",
		NTType: 0,
		Index: 2,
		NumSymbols: 1,
		ReduceFunc: func(X []interface{}) (interface{}, error) {
            return ast.E11(X[0])
		},
	},
	ProdTabEntry{
		String: `T1 : a ;`,
		Id: "T1",
		NTType: 1,
		Index: 3,
		NumSymbols: 1,
		ReduceFunc: func(X []interface{}) (interface{}, error) {
            return ast.T10(X[0])
		},
	},
	
}
