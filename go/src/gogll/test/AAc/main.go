package main

import (
	"gogll/test/AAc/parser/sppf"
	"gogll/test/AAc/parser"
	"gogll/goutil/ioutil"
)

func main() {
	parser.Parse([]byte("c"))
	if err := ioutil.WriteFile("sppf.dot", []byte(sppf.Dot())); err != nil {
		panic(err)
	}
}
