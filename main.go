/*
Copyright 2019 Marius Ackerman

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

package main

import (
	"fmt"
	"os"
	"runtime/pprof"

	"github.com/goccmack/gogll/ast"
	"github.com/goccmack/gogll/frstflw"
	genff "github.com/goccmack/gogll/gen/firstfollow"
	"github.com/goccmack/gogll/gen/golang"
	"github.com/goccmack/gogll/gen/slots"
	gensymbols "github.com/goccmack/gogll/gen/symbols"
	"github.com/goccmack/gogll/gslot"
	"github.com/goccmack/gogll/lexer"
	"github.com/goccmack/gogll/symbols"

	"github.com/goccmack/gogll/cfg"
	"github.com/goccmack/gogll/parser"
	"github.com/goccmack/goutil/md"
)

func main() {
	cfg.GetParams()
	// dumpProcessedMDFile()
	if *cfg.CPUProfile {
		f, err := os.Create("cpu.prof")
		if err != nil {
			fmt.Println("could not create CPU profile: ", err)
			os.Exit(1)
		}
		defer f.Close()
		if err := pprof.StartCPUProfile(f); err != nil {
			fmt.Println("could not start CPU profile: ", err)
			os.Exit(1)
		}
		defer pprof.StopCPUProfile()
	}
	src, err := md.GetSource(cfg.SrcFile)
	if err != nil {
		fail(err)
	}
	t, err := parser.NewParser().Parse(lexer.NewLexer([]byte(src)))
	if err != nil {
		fail(err)
	}
	g := t.(*ast.GoGLL)
	symbols.Init(g)

	gensymbols.Gen(g)
	ff := frstflw.New(g)
	genff.Gen(g, ff)
	gs := gslot.New(g, ff)
	slots.Gen(gs)
	golang.Gen(g, gs, ff)
}

// func dumpProcessedMDFile() {
// 	src, err := md.GetSource(cfg.SrcFile)
// 	if err != nil {
// 		panic(err)
// 	}
// 	ioutil.WriteFile(cfg.SrcFile+".stripped", []byte(src))
// }

func fail(err error) {
	fmt.Printf("Error: %s\n", err)
	os.Exit(1)
}
