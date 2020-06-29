
package parser

import(
	"bytes"
	"fmt"
	"errors"

	parseError "g1/errors"
	"g1/lexer"
	"g1/token"
)

const (
	numProductions 		= 4
	numStates      		= 6
	numTerminals   		= 4
)

// Stack

type stack struct {
	state []int
	attrib	[]interface{}
}

const iNITIAL_STACK_SIZE = 100

func newStack() *stack {
	return &stack{ 	state: 	make([]int, 0, iNITIAL_STACK_SIZE),
					attrib: make([]interface{}, 0, iNITIAL_STACK_SIZE),
			}
}

func (this *stack) reset() {
	this.state = this.state[0:0]
	this.attrib = this.attrib[0:0]
}

func (this *stack) push(s int, a interface{}) {
	this.state = append(this.state, s)
	this.attrib = append(this.attrib, a)
}

func(this *stack) top() int {
	return this.state[len(this.state) - 1]
}

func (this *stack) peek(pos int) int {
	return this.state[pos]
}

func (this *stack) topIndex() int {
	return len(this.state) - 1
}

func (this *stack) popN(items int) []interface{} {
	lo, hi := len(this.state) - items, len(this.state)
	
	attrib := this.attrib[lo: hi]
	
	this.state = this.state[:lo]
	this.attrib = this.attrib[:lo]
	
	return attrib
}

func (this *stack) peekN(items int) []interface{} {
	lo, hi := len(this.state) - items, len(this.state)
	return this.attrib[lo: hi]
}

func (S *stack) String() string {
	w := new(bytes.Buffer)
	fmt.Fprintf(w, "stack:\n")
	for i, st := range S.state {
		fmt.Fprintf(w, "\t%d:%d , ", i, st)
		if S.attrib[i] == nil {
			fmt.Fprintf(w, "nil")
		} else {
			fmt.Fprintf(w, "%v", S.attrib[i])
		}
		w.WriteString("\n")
	}
	return w.String()
}

// Parser

type Parser struct {
	stack     *stack
	nextToken *token.Token

	lex       *lexer.Lexer
	tokens    []*token.Token
	// input position in token stream
	i         int
}

func New(lex *lexer.Lexer) *Parser {
	p := &Parser{
		stack:  newStack(),
		lex:    lex,
		tokens: lex.Tokens,
		i:      0,
	}
	p.stack.push(0, nil)
	return p
}

func (P *Parser) Error(err error) (recovered bool, errorAttrib *parseError.Error) {
	errorAttrib = &parseError.Error{
		Err:            err,
		ErrorToken:     P.nextToken,
		ErrorSymbols:   P.popNonRecoveryStates(),
		ExpectedTokens: make([]string, 0, 8),
	}
	for t, action := range actionTab[P.stack.top()].actions {
		if action != nil {
			errorAttrib.ExpectedTokens = append(errorAttrib.ExpectedTokens, t.ID())
		}
	}

	if action := actionTab[P.stack.top()].actions[token.Error]; action != nil {
		P.stack.push(int(action.(shift)), errorAttrib) // action can only be shift
	} else {
		return
	}

	if action := actionTab[P.stack.top()].actions[P.nextToken.Type()]; action != nil {
		recovered = true
	}
	for !recovered && P.nextToken.Type() != token.EOF {
		P.next()
		if action := actionTab[P.stack.top()].actions[P.nextToken.Type()]; action != nil {
			recovered = true
		}
	}

	return
}

func (P *Parser) popNonRecoveryStates() (removedAttribs []parseError.ErrorSymbol) {
	if rs, ok := P.firstRecoveryState(); ok {
		errorSymbols := P.stack.popN(int(P.stack.topIndex() - rs))
		removedAttribs = make([]parseError.ErrorSymbol, len(errorSymbols))
		for i, e := range errorSymbols {
			removedAttribs[i] = e
		}
	} else {
		removedAttribs = []parseError.ErrorSymbol{}
	}
	return
}

// recoveryState points to the highest state on the stack, which can recover
func (P *Parser) firstRecoveryState() (recoveryState int, canRecover bool) {
	recoveryState, canRecover = P.stack.topIndex(), actionTab[P.stack.top()].canRecover
	for recoveryState > 0 && !canRecover {
		recoveryState--
		canRecover = actionTab[P.stack.peek(recoveryState)].canRecover
	}
	return
}

func (P *Parser) newError(err error) error {
	w := new(bytes.Buffer)
	ln, col := P.nextToken.GetLineColumn()
	fmt.Fprintf(w, "Error @ line %d col %d tok %s", ln, col, P.nextToken)
	if err != nil {
		w.WriteString(err.Error())
	} else {
		w.WriteString(", expected one of: ")
		actRow := actionTab[P.stack.top()]
		for tok, act := range actRow.actions {
			if act != nil {
				fmt.Fprintf(w, "%s ", tok.ID())
			}
		}
	}
	return errors.New(w.String())
}

func (p *Parser) Parse() (res interface{}, err error) {
	p.next()
	for acc := false; !acc; {
		action := actionTab[p.stack.top()].actions[p.nextToken.Type()]

		// fmt.Printf("S%d %s %s\n", p.stack.top(), p.nextToken, action)

		if action == nil {
			if recovered, errAttrib := p.Error(nil); !recovered {
				p.nextToken = errAttrib.ErrorToken
				return nil, p.newError(nil)
			}
			if action = actionTab[p.stack.top()].actions[p.nextToken.Type()]; action == nil {
				panic("Error recovery led to invalid action")
			}
		}

		switch act := action.(type) {
		case accept:
			res = p.stack.popN(1)[0]
			acc = true
		case shift:
			p.stack.push(int(act), p.nextToken)
			p.next()
		case reduce:
			prod := productionsTable[int(act)]
			attrib, err := prod.ReduceFunc(p.stack.popN(prod.NumSymbols))
			if err != nil {
				return nil, p.newError(err)
			} else {
				p.stack.push(gotoTab[p.stack.top()][prod.NTType], attrib)
			}
		default:
			panic("unknown action: " + action.String())
		}
	}
	return res, nil
}

func (p *Parser) next() {
	if p.i < len(p.tokens) {
		p.nextToken = p.tokens[p.i]
		p.i++
	}
}
