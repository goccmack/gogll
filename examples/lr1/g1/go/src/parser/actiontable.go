
package parser

import "g1/token"

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
			token.T_1:shift(3),		/* a */
        },

	},
	actionRow{ // S1
        canRecover: false,
		actions: map[token.Type]action{ 
			token.EOF:accept(true),		/* $ */
			token.T_0:shift(4),		/* + */
        },

	},
	actionRow{ // S2
        canRecover: false,
		actions: map[token.Type]action{ 
			token.EOF:reduce(2),		/* $, reduce: E1 */
			token.T_0:reduce(2),		/* +, reduce: E1 */
        },

	},
	actionRow{ // S3
        canRecover: false,
		actions: map[token.Type]action{ 
			token.EOF:reduce(3),		/* $, reduce: T1 */
			token.T_0:reduce(3),		/* +, reduce: T1 */
        },

	},
	actionRow{ // S4
        canRecover: false,
		actions: map[token.Type]action{ 
			token.T_1:shift(3),		/* a */
        },

	},
	actionRow{ // S5
        canRecover: false,
		actions: map[token.Type]action{ 
			token.EOF:reduce(1),		/* $, reduce: E1 */
			token.T_0:reduce(1),		/* +, reduce: E1 */
        },

	},
}

