package main

import (
	"gogll/test/s0/parser/sppf"
	"gogll/goutil/ioutil"
	"gogll/test/s0/parser"
)

const src = `bac`

func main() {
	parser.Parse([]byte(src))
	if err := ioutil.WriteFile("sppf.dot", []byte(sppf.Dot())); err != nil {
		panic(err)
	}
}