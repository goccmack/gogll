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
	"time"

	"github.com/goccmack/gogll/ast/build"
	"github.com/goccmack/gogll/cfg"
	"github.com/goccmack/gogll/frstflw"
	genff "github.com/goccmack/gogll/gen/firstfollow"
	"github.com/goccmack/gogll/gen/golang"
	"github.com/goccmack/gogll/gen/slots"
	gensymbols "github.com/goccmack/gogll/gen/symbols"
	"github.com/goccmack/gogll/gslot"
	"github.com/goccmack/gogll/lexer"
	"github.com/goccmack/gogll/parser"
	"github.com/goccmack/gogll/parser/bsr"
	"github.com/goccmack/gogll/symbols"
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
	start := time.Now()
	lex := lexer.NewFile(cfg.SrcFile)
	if err, errs := parser.Parse(lex); err != nil {
		fmt.Println(err)
		parseErrors(errs)
	}
	fmt.Printf("Parse duration %s\n", time.Now().Sub(start))

	// bsr.Report()
	g := build.From(bsr.GetRoot(), lex)

	symbols.Init(g)

	gensymbols.Gen(g)
	ff := frstflw.New(g)
	genff.Gen(g, ff)
	gs := gslot.New(g, ff)
	slots.Gen(gs)
	golang.Gen(g, gs, ff)
}

func fail(err error) {
	fmt.Printf("Error: %s\n", err)
	os.Exit(1)
}

func parseErrors(errs []*parser.Error) {
	fmt.Println("Parse Errors:")
	for _, err := range errs {
		fmt.Println(err)
	}
	os.Exit(1)
}
