// Package lexfsa generates a text file containing the lexer FSA
package lexfsa

import (
	"bytes"
	"fmt"

	"github.com/goccmack/goutil/ioutil"

	"github.com/goccmack/gogll/lex/items"
)

func Gen(fname string, ls *items.Sets) {
	w := new(bytes.Buffer)
	for _, s := range ls.Sets() {
		fmt.Fprintf(w, "S%d:\n", s.No)
		for _, i := range s.Items() {
			fmt.Fprintf(w, "    %s\n", i)
		}

		fmt.Fprintln(w, "  Transitions:")
		for _, e := range s.Transitions {
			fmt.Fprintf(w, "    %s -> S%d\n", e.Event, e.To.No)
		}

		fmt.Fprintln(w)
	}
	if err := ioutil.WriteFile(fname, w.Bytes()); err != nil {
		panic(err)
	}
}
