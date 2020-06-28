
package parser

import(
    "logs/ast"
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
		String: `G0 : Lines ;`,
		Id: "G0",
		NTType: 0,
		Index: 0,
		NumSymbols: 1,
		ReduceFunc: func(X []interface{}) (interface{}, error) {
            return ast.G00(X[0])
		},
	},
	ProdTabEntry{
		String: `Lines : Line ;`,
		Id: "Lines",
		NTType: 1,
		Index: 1,
		NumSymbols: 1,
		ReduceFunc: func(X []interface{}) (interface{}, error) {
            return ast.Lines0(X[0])
		},
	},
	ProdTabEntry{
		String: `Lines : Lines Line ;`,
		Id: "Lines",
		NTType: 1,
		Index: 2,
		NumSymbols: 2,
		ReduceFunc: func(X []interface{}) (interface{}, error) {
            return ast.Lines1(X[0],X[1])
		},
	},
	ProdTabEntry{
		String: `Line : sap ip name name timestamp string number1 number1 string string ;`,
		Id: "Line",
		NTType: 0,
		Index: 3,
		NumSymbols: 10,
		ReduceFunc: func(X []interface{}) (interface{}, error) {
            return ast.Line0(X[0],X[1],X[2],X[3],X[4],X[5],X[6],X[7],X[8],X[9])
		},
	},
	
}
