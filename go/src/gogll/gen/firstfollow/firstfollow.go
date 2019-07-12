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

func Gen() {
	w := new(bytes.Buffer)
	genFirstSets(w)
	genFollowSets(w)
	fname := path.Join(cfg.BaseDir, "first_follow.txt")
	if err := ioutil.WriteFile(fname, w.Bytes()); err != nil {
		fmt.Printf("Error writing first and follow: %s\n", err)
		os.Exit(1)
	}
}

func genFirstSets(w io.Writer) {
	for _, s := range ast.GetSymbols() {
		genFirstSet(w, s)
	}
}

func genFirstSet(w io.Writer, symbol string) {
	fmt.Fprintf(w, "%s: ", symbol)
	for _, s := range frstflw.FirstOfSymbol(symbol).Elements() {
		fmt.Fprintf(w, "%s ", s)
	}
	fmt.Fprintln(w)
}

func genFollowSets(w io.Writer) {
	for _, s := range ast.GetNonTerminals() {
		genFollowSet(w, s)
	}
}

func genFollowSet(w io.Writer, nt string) {
	fmt.Fprintf(w, "Follow(%s): ", nt)
	for _, s := range frstflw.Follow(nt).Elements() {
		fmt.Fprintf(w, "%s ", s)
	}
	fmt.Fprintln(w)
}
