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

package main

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"runtime/pprof"

	"github.com/goccmack/gogll/ast"
	"github.com/goccmack/gogll/cfg"
	"github.com/goccmack/gogll/frstflw"
	genff "github.com/goccmack/gogll/gen/firstfollow"
	gengogll "github.com/goccmack/gogll/gen/golang/gll"
	gengolexer "github.com/goccmack/gogll/gen/golang/lexer"
	gengolr1 "github.com/goccmack/gogll/gen/golang/lr1"
	gengotoken "github.com/goccmack/gogll/gen/golang/token"
	"github.com/goccmack/gogll/gen/lexfsa"
	genrustgll "github.com/goccmack/gogll/gen/rust/gll"
	genrustlexer "github.com/goccmack/gogll/gen/rust/lexer"
	genrustlr1 "github.com/goccmack/gogll/gen/rust/lr1"
	genrusttoken "github.com/goccmack/gogll/gen/rust/token"
	"github.com/goccmack/gogll/gen/slots"
	gensymbols "github.com/goccmack/gogll/gen/symbols"
	"github.com/goccmack/gogll/gslot"
	"github.com/goccmack/gogll/im/tokens"
	"github.com/goccmack/gogll/lex/items"
	"github.com/goccmack/gogll/lexer"
	"github.com/goccmack/gogll/lr1"
	"github.com/goccmack/gogll/parser"
	"github.com/goccmack/gogll/sc"
	"github.com/goccmack/gogll/symbols"
)

func main() {
	cfg.GetParams()
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
	lex := lexer.NewFile(cfg.SrcFile)
	bsrSet, errs := parser.Parse(lex)
	if errs != nil {
		parseErrors(errs)
	}

	if bsrSet.IsAmbiguous() {
		fmt.Println("Error: Ambiguous parse forest")
		bsrSet.ReportAmbiguous()
		os.Exit(1)
	}

	g := ast.Build(bsrSet.GetRoot(), lex)
	sc.Go(g, lex)
	symbols.Init(g)

	ff := frstflw.New(g)
	gs := gslot.New(g, ff)
	ts := tokens.New(g)

	lexSets := items.New(g)
	if cfg.Verbose {
		gensymbols.Gen(g)
		genff.Gen(g, ff)
		slots.Gen(gs)
		lexfsa.Gen(filepath.Join(cfg.BaseDir, "lexfsa.txt"), lexSets)
	}

	switch cfg.Target {
	case cfg.Go:
		gengolexer.Gen(g, lexSets, ts)
		gengotoken.Gen(g, ts)
		if len(g.SyntaxRules) > 0 {
			if *cfg.GLL {
				gengogll.Gen(g, gs, ff, ts)
			} else {
				bprods, states, actions := lr1.Gen(g)
				gengolr1.Gen(g.Package.GetString(), bprods, states, actions)
			}
		}
	case cfg.Rust:
		genrusttoken.Gen(filepath.Join(cfg.BaseDir, "src", "token", "mod.rs"), ts)
		genrustlexer.Gen(path.Join(cfg.BaseDir, "src", "lexer", "mod.rs"), g, lexSets, ts)
		if len(g.SyntaxRules) > 0 {
			if *cfg.GLL {
				genrustgll.Gen(path.Join(cfg.BaseDir, "src", "parser"), g, gs, ff, ts)
			} else {
				bprods, states, actions := lr1.Gen(g)
				genrustlr1.Gen(g.Package.GetString(), bprods, states, actions)
			}
		}
	default:
		fail(fmt.Errorf("Invalid target language"))
	}

}

func fail(err error) {
	fmt.Printf("Error: %s\n", err)
	os.Exit(1)
}

func parseErrors(errs []*parser.Error) {
	fmt.Println("Parse Errors:")
	ln := errs[0].Line
	for _, err := range errs {
		if err.Line == ln {
			fmt.Println(err)
		}
	}
	os.Exit(1)
}
