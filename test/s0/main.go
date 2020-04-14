package main

import (
	"github.com/goccmack/gogll/test/s0/parser"
	"github.com/goccmack/gogll/test/s0/parser/sppf"
	"github.com/goccmack/goutil/ioutil"
)

const src = `bac`

func main() {
	parser.Parse([]byte(src))
	if err := ioutil.WriteFile("sppf.dot", []byte(sppf.Dot())); err != nil {
		panic(err)
	}
}
