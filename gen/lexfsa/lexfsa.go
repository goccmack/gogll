/*
Copyright 2020 Marius Ackerman

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

  http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Package lexfsa generates a text file containing the lexer FSA
package lexfsa

import (
	"bytes"
	"fmt"

	"github.com/goccmack/goutil/ioutil"

	"github.com/goccmack/gogll/v3/lex/items"
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
