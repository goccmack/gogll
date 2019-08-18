package main

import (
	"fmt"
	"gogll/cfg"
	"gogll/goutil/bsr"
	"gogll/parser"
)

func main() {
	cfg.GetParams()
	if err, errs := parser.ParseFile(cfg.SrcFile); err != nil {
		fail(err, errs)
	}
}

func fail(err error, errs []*parser.ParseError) {
	fmt.Printf("ParseError: %s\n", err)
	// parser.DumpCRF(errs[0].InputPos)
	bsr.Dump()
	for _, e := range errs {
		fmt.Println("", e)
	}
}
