package cfg

import (
	"flag"
	"fmt"
	"os"
	"path"
)

var (
	BaseDir string
	BNFFile string
	Package string

	pkg = flag.String("p", "", "")
)

func GetParams() {
	flag.Parse()
	getSourceFile()
	getFileBase()
	getPackage()
}

func getFileBase() {
	BaseDir, _ = path.Split(BNFFile)
	if BaseDir == "" {
		BaseDir = "."
	}
}

func getPackage() {
	if *pkg == "" {
		fail("Package prefix must be specified")
	}
	Package = *pkg
}

func getSourceFile() {
	if flag.NArg() < 1 {
		fail("Source file required")
	}
	BNFFile = flag.Arg(0)
}

func fail(msg string) {
	fmt.Printf("ERROR: %s\n", msg)
	usage()
	os.Exit(1)
}

func usage() {
	msg := `use: gogll -p <package> <bnfFile file>
	<file name> : Mandatory. Name of the BNF file to be processed
    -p <package>: Mandatory. Package prefix of the generated parser packages
                  E.g.: -p test generates parser files with a package: "test/parser"`
	fmt.Println(msg)
}
