
/*
*/
package parser

const numNTSymbols = 2
type(
	gotoTable [numStates]gotoRow
	gotoRow	[numNTSymbols] int
)

var gotoTab = gotoTable{
	gotoRow{ // S0
		1, // E1
        2, // T1
        
	},
	gotoRow{ // S1
		-1, // E1
        -1, // T1
        
	},
	gotoRow{ // S2
		-1, // E1
        -1, // T1
        
	},
	gotoRow{ // S3
		-1, // E1
        -1, // T1
        
	},
	gotoRow{ // S4
		-1, // E1
        5, // T1
        
	},
	gotoRow{ // S5
		-1, // E1
        -1, // T1
        
	},
	
}
