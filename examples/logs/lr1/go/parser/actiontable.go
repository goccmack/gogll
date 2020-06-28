
package parser

import "logs/token"

type(
    actionTable [numStates]actionRow
    actionRow struct {
        canRecover bool
        actions map[token.Type]action
    }
)

var actionTab = actionTable{ 
	actionRow{ // S0
        canRecover: false,
		actions: map[token.Type]action{ 
			token.T_3:shift(3),		/* sap */
        },

	},
	actionRow{ // S1
        canRecover: false,
		actions: map[token.Type]action{ 
			token.EOF:reduce(1),		/* $, reduce: Lines */
			token.T_3:reduce(1),		/* sap, reduce: Lines */
        },

	},
	actionRow{ // S2
        canRecover: false,
		actions: map[token.Type]action{ 
			token.EOF:accept(true),		/* $ */
			token.T_3:shift(3),		/* sap */
        },

	},
	actionRow{ // S3
        canRecover: false,
		actions: map[token.Type]action{ 
			token.T_0:shift(5),		/* ip */
        },

	},
	actionRow{ // S4
        canRecover: false,
		actions: map[token.Type]action{ 
			token.EOF:reduce(2),		/* $, reduce: Lines */
			token.T_3:reduce(2),		/* sap, reduce: Lines */
        },

	},
	actionRow{ // S5
        canRecover: false,
		actions: map[token.Type]action{ 
			token.T_1:shift(6),		/* name */
        },

	},
	actionRow{ // S6
        canRecover: false,
		actions: map[token.Type]action{ 
			token.T_1:shift(7),		/* name */
        },

	},
	actionRow{ // S7
        canRecover: false,
		actions: map[token.Type]action{ 
			token.T_5:shift(8),		/* timestamp */
        },

	},
	actionRow{ // S8
        canRecover: false,
		actions: map[token.Type]action{ 
			token.T_4:shift(9),		/* string */
        },

	},
	actionRow{ // S9
        canRecover: false,
		actions: map[token.Type]action{ 
			token.T_2:shift(10),		/* number1 */
        },

	},
	actionRow{ // S10
        canRecover: false,
		actions: map[token.Type]action{ 
			token.T_2:shift(11),		/* number1 */
        },

	},
	actionRow{ // S11
        canRecover: false,
		actions: map[token.Type]action{ 
			token.T_4:shift(12),		/* string */
        },

	},
	actionRow{ // S12
        canRecover: false,
		actions: map[token.Type]action{ 
			token.T_4:shift(13),		/* string */
        },

	},
	actionRow{ // S13
        canRecover: false,
		actions: map[token.Type]action{ 
			token.EOF:reduce(3),		/* $, reduce: Line */
			token.T_3:reduce(3),		/* sap, reduce: Line */
        },

	},
}

