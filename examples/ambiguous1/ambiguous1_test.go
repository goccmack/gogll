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

package ambiguous1

import (
	"testing"

	"github.com/goccmack/gogll/examples/ambiguous1/goutil/bsr"
	"github.com/goccmack/gogll/examples/ambiguous1/parser"
	"github.com/goccmack/gogll/examples/ambiguous1/parser/symbols"
)

func Test1(t *testing.T) {
	parser.Parse([]byte("aba"))
	da()
	bsr.DumpSlots()
}

func da() {
	rmBa()
	rmBereavedParents()
}

func rmBa() {
	for _, b := range bsr.GetAll() {
		// fmt.Println(b)
		if b.GetString() == "a" && b.Label.Head() == "B" {
			// fmt.Println("  ignore")
			b.Ignore()
		}
	}
}

func rmBereavedParents() {
	for again := true; again; {
		again = false
		for _, b := range bsr.GetAll() {
			for i, c := range b.Label.Symbols() {
				if symbols.IsNonTerminal(c) && nil == b.GetNTChildrenI(i) {
					b.Ignore()
					again = true
				}
			}
		}
	}
}
