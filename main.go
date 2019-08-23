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
	"gogll/cfg"
	"gogll/da"
	"gogll/frstflw"
	genff "gogll/gen/firstfollow"
	"gogll/gen/golang"
	"gogll/gen/slots"
	"gogll/gen/symbols"
	"gogll/goutil/bsr"
	"gogll/gslot"
	"gogll/parser"
	"gogll/sa"
	"io/ioutil"
	"os"
)

func main() {
	cfg.GetParams()
	if err, errs := parser.ParseFile(cfg.SrcFile); err != nil {
		fail(err, errs)
	}
	da.Go()
	g, errs := sa.Go()
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
}

func fail(err error, errs []*parser.ParseError) {
	fmt.Printf("ParseError: %s\n", err)
	// parser.DumpCRF(errs[0].InputPos)
	bsr.Dump()
	for _, e := range errs {
		fmt.Println("", e)
	}
}

func getInput() string {
	buf, err := ioutil.ReadFile(cfg.SrcFile)
	if err != nil {
		panic(err)
	}
	return string(buf)
}
