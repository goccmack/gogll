//  Copyright 2019 Marius Ackerman
//
//  Licensed under the Apache License, Version 2.0 (the "License");
//  you may not use this file except in compliance with the License.
//  You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
//  Unless required by applicable law or agreed to in writing, software
//  distributed under the License is distributed on an "AS IS" BASIS,
//  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//  See the License for the specific language governing permissions and
//  limitations under the License.
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
	BaseDir, _ = path.Split(SrcFile)
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
	SrcFile = flag.Arg(0)
}

func fail(msg string) {
	fmt.Printf("ERROR: %s\n", msg)
	usage()
	os.Exit(1)
}

func usage() {
	msg := `use: gogll -p <package> <source file>
	<source file> : Mandatory. Name of the source file to be processed. 
					If the file extension is ".md" the bnf is extracted from markdown code segments
                    enclosed in triple backticks.
    -p <package>  : Mandatory. Package prefix of the generated parser packages
                    E.g.: -p test generates parser files with a package: "test/parser"`
	fmt.Println(msg)
}
