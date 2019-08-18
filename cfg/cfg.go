package cfg

import (
	"flag"
	"fmt"
	"os"
	"path"
)

var (
	BaseDir string
	SrcFile string
)

func GetParams() {
	flag.Parse()
	getSourceFile()
	getFileBase()
}

func getFileBase() {
	BaseDir, _ = path.Split(SrcFile)
	if BaseDir == "" {
		BaseDir = "."
	}
}

func getSourceFile() {
	if flag.NArg() < 1 {
		fail("Source file required")
	}
	SrcFile = flag.Arg(0)
}

func fail(msg string) {
	fmt.Printf("ERROR: %s\n", msg)
	usage()
	os.Exit(1)
}

func usage() {
	msg := `use: gogll <source file>
	<source file> : Mandatory. Name of the source file to be processed. 
					If the file extension is ".md" the bnf is extracted from markdown code segments
                    enclosed in triple backticks.`
	fmt.Println(msg)
}
