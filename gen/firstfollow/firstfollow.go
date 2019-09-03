//  Copyright 2019 Marius Ackerman
//
//  Licensed under the Apache License, Version 2.0 (the "License");
//  you may not use this file except in compliance with the License.
//  You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
//  Unless required by applicable law or agreed to in writing, software
//  distributed under the License is distributed on an "AS IS" BASIS,
//  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//  See the License for the specific language governing permissions and
//  limitations under the License.

package firstfollow

import (
	"bytes"
	"fmt"
	"github.com/goccmack/gogll/ast"
	"github.com/goccmack/gogll/cfg"
	"github.com/goccmack/gogll/frstflw"
	"github.com/goccmack/gogll/goutil/ioutil"
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
