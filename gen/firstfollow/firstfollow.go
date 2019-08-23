package firstfollow

import (
	"bytes"
	"fmt"
	"gogll/ast"
	"gogll/cfg"
	"gogll/frstflw"
	"gogll/goutil/ioutil"
	"io"
	"os"
	"path"
)

func Gen(g *ast.Grammar, ff *frstflw.FF) {
	w := new(bytes.Buffer)
	genFirstSets(w, g, ff)
	genFollowSets(w, g, ff)
	fname := path.Join(cfg.BaseDir, "first_follow.txt")
	if err := ioutil.WriteFile(fname, w.Bytes()); err != nil {
		fmt.Printf("Error writing first and follow: %s\n", err)
		os.Exit(1)
	}
}

func genFirstSets(w io.Writer, g *ast.Grammar, ff *frstflw.FF) {
	for _, s := range g.GetSymbols() {
		genFirstSet(w, s, ff)
	}
}

func genFirstSet(w io.Writer, symbol string, ff *frstflw.FF) {
	fmt.Fprintf(w, "%s: ", symbol)
	for _, s := range ff.FirstOfSymbol(symbol).Elements() {
		fmt.Fprintf(w, `%s `, s)
	}
	fmt.Fprintln(w)
}

func genFollowSets(w io.Writer, g *ast.Grammar, ff *frstflw.FF) {
	for _, s := range g.GetNonTerminals() {
		genFollowSet(w, s, ff)
	}
}

func genFollowSet(w io.Writer, nt string, ff *frstflw.FF) {
	fmt.Fprintf(w, "Follow(%s): ", nt)
	for _, s := range ff.Follow(nt).Elements() {
		fmt.Fprintf(w, `%s `, s)
	}
	fmt.Fprintln(w)
}
