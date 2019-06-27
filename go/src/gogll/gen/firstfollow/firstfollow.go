package firstfollow

import (
	"bytes"
	"fmt"
	"gogll/ast"
	"io"
	"io/ioutil"
	"os"
	"path"
)

func Gen(dir string, g *ast.Grammar) {
	w := new(bytes.Buffer)
	genFirstSets(w, g)
	genFollowSets(w, g)
	fname := path.Join(dir, "first_follow.txt")
	if err := ioutil.WriteFile(fname, w.Bytes(), 0731); err != nil {
		fmt.Printf("Error writing first and follow: %s\n", err)
		os.Exit(1)
	}
}

func genFirstSets(w io.Writer, g *ast.Grammar) {
	for _, s := range ast.GetSymbols() {
		genFirstSet(w, g, s)
	}
}

func genFirstSet(w io.Writer, g *ast.Grammar, symbol string) {
	fmt.Fprintf(w, "%s: ", symbol)
	for _, s := range g.FirstOfSymbol(symbol).Elements() {
		fmt.Fprintf(w, "%s ", s)
	}
	fmt.Fprintln(w)
}

func genFollowSets(w io.Writer, g *ast.Grammar) {
	for _, s := range ast.GetNonTerminals() {
		genFollowSet(w, g, s)
	}
}

func genFollowSet(w io.Writer, g *ast.Grammar, nt string) {
	fmt.Fprintf(w, "Follow(%s): ", nt)
	for _, s := range g.Follow(nt).Elements() {
		fmt.Fprintf(w, "%s ", s)
	}
	fmt.Fprintln(w)
}
