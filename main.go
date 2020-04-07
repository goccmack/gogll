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

	"github.com/goccmack/gogll/goutil/ioutil"

	"github.com/goccmack/gogll/cfg"
	"github.com/goccmack/gogll/da"
	"github.com/goccmack/gogll/frstflw"
	genff "github.com/goccmack/gogll/gen/firstfollow"
	"github.com/goccmack/gogll/gen/golang"
	"github.com/goccmack/gogll/gen/slots"
	"github.com/goccmack/gogll/gen/symbols"
	"github.com/goccmack/gogll/goutil/md"
	"github.com/goccmack/gogll/gslot"
	"github.com/goccmack/gogll/parser"
	"github.com/goccmack/gogll/sa"
)

var (
	parseDur time.Duration
	daDur    time.Duration
	saDur    time.Duration
	genDur   time.Duration
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
	startTime := time.Now()
	if err, errs := parser.ParseFile(cfg.SrcFile); err != nil {
		fail(err, errs)
	}
	parseDur = time.Now().Sub(startTime)

	// da.Report()

	startTime = time.Now()
	da.Go()
	daDur = time.Now().Sub(startTime)

	// fmt.Println("Ambiguous BSRs after disambiguation")
	// da.Report()

	startTime = time.Now()
	g, errs := sa.Go()
	saDur = time.Now().Sub(startTime)
	startTime = time.Now()
	if errs != nil {
		for _, err := range errs {
			fmt.Println(err)
		}
		os.Exit(1)
	}
	symbols.Gen(g)
	ff := frstflw.New(g)
	genff.Gen(g, ff)
	gs := gslot.New(g, ff)
	slots.Gen(gs)
	golang.Gen(g, gs, ff)
	genDur = time.Now().Sub(startTime)
	fmt.Printf("parse %.3f ms, da %.3f ms, sa %.3f ms, gen %.3f ms\n",
		float64(parseDur)/float64(time.Millisecond),
		float64(daDur)/float64(time.Millisecond),
		float64(saDur)/float64(time.Millisecond),
		float64(genDur)/float64(time.Millisecond))
	// if *cfg.BSRStats {
	// 	for r, c := range bsr.Stats() {
	// 		fmt.Println(r, c)
	// 	}
	// }
}

func dumpProcessedMDFile() {
	src, err := md.GetSource(cfg.SrcFile)
	if err != nil {
		panic(err)
	}
	ioutil.WriteFile(cfg.SrcFile+".stripped", []byte(src))
}

func fail(err error, errs []*parser.ParseError) {
	fmt.Printf("ParseError: %s\n", err)
	// parser.DumpCRF(errs[0].InputPos)
	// bsr.Dump()
	for _, e := range errs {
		fmt.Println("", e)
	}
}
