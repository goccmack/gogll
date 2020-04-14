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
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/goccmack/goutil/ioutil"

	"github.com/goccmack/gogll/examples/ambiguous1/goutil/bsr"
	"github.com/goccmack/gogll/examples/ambiguous1/parser"
)

type disambiguator int

type dotWriter struct {
	fname       string
	w           *bytes.Buffer
	nodeNum     int
	nodeNameMap map[string]string // BSR.GetString() -> string
	linkMap     map[link]bool
}

type link struct {
	from, to string
}

type walker interface {
	do(node, parent bsr.BSR)
}

var nilBSR = bsr.BSR{}

func Test1(t *testing.T) {
	parser.Parse([]byte("aba"))

	dw := newDotWriter("beforeDA.dot")
	for _, r := range bsr.GetRoots() {
		walk(r, dw, nilBSR)
	}
	dw.close()

	for _, r := range bsr.GetRoots() {
		walk(r, disambiguator(0), nilBSR)
	}

	dw = newDotWriter("afterDA.dot")
	for _, r := range bsr.GetRoots() {
		walk(r, dw, nilBSR)
	}
	dw.close()
}

/*** BSR walker ***/

func walk(node bsr.BSR, w walker, parent bsr.BSR) {
	w.do(node, parent)
	switch node.Label.Head() {
	case "S":
		switch node.Alternate() {
		case 0: // S : A S
			for _, a := range node.GetNTChildren("A", 0) {
				walk(a, w, node)
			}
			for _, s := range node.GetNTChildren("S", 0) {
				walk(s, w, node)
			}
		case 1: // S : B S
			for _, b1 := range node.GetNTChildren("B", 0) {
				walk(b1, w, node)
			}
			for _, s := range node.GetNTChildren("S", 0) {
				walk(s, w, node)
			}
		case 2: // S : empty
			// ignore
		}
	case "A":
		// ignore
	case "B":
		// ignore
	}
}

/*** Disabiguate BSR set ***/

func (da disambiguator) do(s, parent bsr.BSR) {
	switch s.Label.Head() {
	case "S":
		daS(s)
	case "A":
		// nothing required
	case "B":
		daB(s)
	}
}

func daS(s bsr.BSR) {
	switch s.Alternate() {
	case 0: // S : A S
		for _, s1 := range s.GetNTChildren("S", 0) {
			daS(s1)
		}
		if len(s.GetNTChildren("S", 0)) == 0 {
			s.Ignore()
		}
	case 1: // S : B S
		for _, b := range s.GetNTChildren("B", 0) {
			daB(b)
		}
		if len(s.GetNTChildren("B", 0)) == 0 {
			s.Ignore()
		}
		for _, s1 := range s.GetNTChildren("S", 0) {
			daS(s1)
		}
		if len(s.GetNTChildren("S", 0)) == 0 {
			s.Ignore()
		}
	case 2: // S : empty
		// do nothing
	}
}

func daB(b bsr.BSR) {
	if b.GetString() == "a" {
		b.Ignore()
	}
}

/*** Write dot files***/

func newDotWriter(fname string) *dotWriter {
	dw := &dotWriter{
		fname:       fname,
		w:           new(bytes.Buffer),
		nodeNameMap: make(map[string]string),
		linkMap:     make(map[link]bool),
	}
	fmt.Fprintf(dw.w, "digraph %s {\n", strings.Split(fname, ".")[0])
	return dw
}

func (dw *dotWriter) newNodeName(b bsr.BSR) string {
	if nnm, exist := dw.nodeNameMap[b.String()]; exist {
		return nnm
	}

	nodeName := fmt.Sprintf("Node%d", dw.nodeNum)
	dw.nodeNum++
	dw.nodeNameMap[b.String()] = nodeName
	return nodeName
}

func (dw *dotWriter) do(s, parent bsr.BSR) {
	nnm := dw.newNodeName(s)
	fmt.Fprintf(dw.w, "%s [label=\"%s\"]\n", nnm, s)
	if parent != nilBSR {
		lnk := link{parent.String(), s.String()}
		if _, exist := dw.linkMap[lnk]; !exist {
			fmt.Fprintf(dw.w, "%s -> %s\n", dw.nodeNameMap[parent.String()], nnm)
			dw.linkMap[lnk] = true
		}
	}
}

func (dw *dotWriter) close() {
	fmt.Fprintln(dw.w, "}")
	ioutil.WriteFile(dw.fname, dw.w.Bytes())
}
