/*
Package lr1 generates Knuth's orgininal LR(1) parser and Pager's PGM with
weak compatibility.
*/
package lr1

import (
	"bytes"
	"fmt"
	"os"
	"path"

	"github.com/goccmack/gogll/ast"
	"github.com/goccmack/gogll/cfg"
	"github.com/goccmack/gogll/lr1/action"
	"github.com/goccmack/gogll/lr1/basicprod"
	"github.com/goccmack/gogll/lr1/first"
	"github.com/goccmack/gogll/lr1/items"
	"github.com/goccmack/gogll/lr1/knuth"
	"github.com/goccmack/gogll/lr1/pgm"
	"github.com/goccmack/gogll/lr1/states"
	"github.com/goccmack/gogll/symbols"
	"github.com/goccmack/goutil/ioutil"
)

func Gen(g *ast.GoGLL) ([]*basicprod.Production, *states.States, action.Actions) {
	removeOldFiles()

	prods := basicprod.Get(g.SyntaxRules)
	items := items.NewItems(prods)
	smbls := symbols.GetSymbols()
	first := first.New(prods)
	var states *states.States
	if *cfg.Knuth {
		//TODO: remove symbols
		states = knuth.States(smbls, items, first)
	} else {
		//TODO: remove symbols
		states = pgm.States(smbls, items, first)
	}

	actions, conflicts := action.GetActions(states)
	handleConflicts(conflicts, states)

	// if cfg.Verbose2() {
	// writeBasicProductions(cfg, prods)
	ioutil.WriteFile(path.Join(cfg.BaseDir, "CFG_items.txt"), []byte(items.String()))
	// writeCFGSymbols(symbols)
	// io.WriteFileString(path.Join(cfg.BaseDir, "first.txt"), first.String())
	ioutil.WriteFile(path.Join(cfg.BaseDir, "LR1_states.txt"), statesString(states, actions))
	// ioutil.WriteFile(path.Join(cfg.BaseDir, "LR_states.dot"), statesDotString(states, actions, prods))
	// }

	return prods, states, actions
}

func handleConflicts(conflicts [][]*action.Conflict, states *states.States) {
	if numConflicts(conflicts) > 0 {
		fmt.Printf("%d LR(1) conflicts. See LR1_conflicts.txt", numConflicts(conflicts))
		writeConflicts(conflicts)
		if !*cfg.AutoResolveLRConf {
			fmt.Println("LR(1) conflicts were automatically resolved")
			os.Exit(1)
		}
	}
}

func numConflicts(conflicts [][]*action.Conflict) (num int) {
	for _, sc := range conflicts {
		for _, c := range sc {
			if c != nil {
				num++
			}
		}
	}
	return
}

func writeConflicts(conflicts [][]*action.Conflict) {
	w := new(bytes.Buffer)
	fmt.Fprintf(w, "%d LR(1) Conflicts:\n", numConflicts(conflicts))
	conflictNo := 1
	for si, sc := range conflicts {
		for _, c := range sc {
			fmt.Fprintf(w, "%4d) S%d: %s\n", conflictNo, si, c)
			conflictNo++
		}
	}
	ioutil.WriteFile(path.Join(cfg.BaseDir, "LR1_conflicts.txt"), w.Bytes())
}

// func reduceLabel(state *states.State, actions map[string]action.Action, prods []*basicprod.Production) string {
// 	w := new(bytes.Buffer)
// 	numRedux, accept := 0, false
// 	for sym, a := range actions {
// 		switch red := a.(type) {
// 		case action.Accept:
// 			accept = true
// 		case action.Reduce:
// 			if numRedux > 0 {
// 				fmt.Fprintf(w, "\\n")
// 			}
// 			fmt.Fprintf(w, "%s {%s}", prods[red].DotString(), sym)
// 			numRedux++
// 		}
// 	}
// 	switch {
// 	case numRedux > 0:
// 		return fmt.Sprintf("S%d [shape=box,label=\"S%d\\n%s\"]\n", state.Number, state.Number, w.String())
// 	case accept:
// 		return fmt.Sprintf("S%d [shape=doublecircle]\n", state.Number)
// 	default:
// 		return ""
// 	}
// }

// func statesDotString(states *states.States, actions action.Actions, prods []*basicprod.Production) string {
// 	w := new(bytes.Buffer)
// 	fmt.Fprintf(w, "digraph{\n")
// 	for si, state := range states.List {
// 		fmt.Fprintf(w, "%s", reduceLabel(state, actions[si], prods))
// 		for _, t := range state.Transitions.List() {
// 			fmt.Fprintf(w, "\tS%d ->S%d [label = \"%s\"]\n", state.Number, t.State.Number, t.Sym)
// 		}
// 	}
// 	fmt.Fprintf(w, "}\n")
// 	return w.String()
// }

func statesString(states *states.States, actions action.Actions) []byte {
	w := new(bytes.Buffer)
	for si, state := range states.List {
		fmt.Fprintf(w, "%s", state)
		fmt.Fprintf(w, "Actions:\n")
		if actions != nil {
			for _, symT := range symbols.GetTerminalSymbols() {
				if action := actions[si][symT]; action != nil {
					fmt.Fprintf(w, "\t%s: %s\n", symT, action)
				}
			}
		}
		fmt.Fprintln(w)
	}
	return w.Bytes()
}

func removeOldFiles() {
	os.Remove(path.Join(cfg.BaseDir, "basic_productions.txt"))
	os.Remove(path.Join(cfg.BaseDir, "CFG_items.txt"))
	os.Remove(path.Join(cfg.BaseDir, "CFG_symbols.txt"))
	os.Remove(path.Join(cfg.BaseDir, "first.txt"))
	os.Remove(path.Join(cfg.BaseDir, "LR1_states.txt"))
	os.Remove(path.Join(cfg.BaseDir, "LR1_conflicts.txt"))
}

// func writeBasicProductions(cfg config.Config, prods []*basicprod.Production) {
// 	w := new(bytes.Buffer)
// 	for i, prod := range prods {
// 		fmt.Fprintf(w, "%4d: %s\n\n", i, prod)
// 	}
// 	io.WriteFileString(path.Join(cfg.BaseDir, "basic_productions.txt"), w.String())
// }

// func writeCFGSymbols(cfg config.Config, symbols symbols.Symbols) {
// 	w := new(bytes.Buffer)
// 	fmt.Fprintf(w, "Start Symbol (S): %s\n\n", symbols.StartSymbol)
// 	fmt.Fprintf(w, "Terminal Symbols (T):\n")
// 	for _, t := range symbols.ListTerminals() {
// 		fmt.Fprintf(w, "\t%s\n", t)
// 	}
// 	fmt.Fprintf(w, "\n")
// 	fmt.Fprintf(w, "Non-terminal Symbols (NT):\n")
// 	for _, nt := range symbols.ListNonTerminals() {
// 		fmt.Fprintf(w, "\t%s\n", nt)
// 	}
// 	io.WriteFileString(path.Join(cfg.BaseDir, "CFG_symbols.txt"), w.String())
// }
