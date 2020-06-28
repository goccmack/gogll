
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
		1, // Line
        2, // Lines
        
	},
	gotoRow{ // S1
		-1, // Line
        -1, // Lines
        
	},
	gotoRow{ // S2
		4, // Line
        -1, // Lines
        
	},
	gotoRow{ // S3
		-1, // Line
        -1, // Lines
        
	},
	gotoRow{ // S4
		-1, // Line
        -1, // Lines
        
	},
	gotoRow{ // S5
		-1, // Line
        -1, // Lines
        
	},
	gotoRow{ // S6
		-1, // Line
        -1, // Lines
        
	},
	gotoRow{ // S7
		-1, // Line
        -1, // Lines
        
	},
	gotoRow{ // S8
		-1, // Line
        -1, // Lines
        
	},
	gotoRow{ // S9
		-1, // Line
        -1, // Lines
        
	},
	gotoRow{ // S10
		-1, // Line
        -1, // Lines
        
	},
	gotoRow{ // S11
		-1, // Line
        -1, // Lines
        
	},
	gotoRow{ // S12
		-1, // Line
        -1, // Lines
        
	},
	gotoRow{ // S13
		-1, // Line
        -1, // Lines
        
	},
	
}
