//  Copyright 2020 Marius Ackerman
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

// Version is the version of this compiler
const Version = "v3.0.2"

var (
	BaseDir    string
	SrcFile    string
	Verbose    bool
	BSRStats   = flag.Bool("bs", false, "Print BSR stats")
	help       = flag.Bool("h", false, "Print help")
	CPUProfile = flag.Bool("CPUProf", false, "Generate CPU profile")
	outDir     = flag.String("o", "", "")
	verbose    = flag.Bool("v", false, "Verbose")
	version    = flag.Bool("version", false, "Version")
)

func GetParams() {
	flag.Parse()
	if *help {
		usage()
		os.Exit(0)
	}
	if *version {
		fmt.Println("gogll", Version)
		os.Exit(0)
	}
	getSourceFile()
	getFileBase()
	Verbose = *verbose
}

func getFileBase() {
	if *outDir != "" {
		BaseDir = *outDir
	} else {
		BaseDir, _ = path.Split(SrcFile)
		if BaseDir == "" {
			BaseDir = "."
		}
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
	msg := `use: gogll [-h][-version][-v][-CPUProf] [-o <out dir>] <source file>
    
    -o <out dir>: Optional. The directory to which code will be generated.
                  Default: the same directory as <source file>.
    
    <source file> : Mandatory. Name of the source file to be processed. 
        If the file extension is ".md" the bnf is extracted from markdown code 
        segments enclosed in triple backticks.

    -CPUProf : Optional. Generate a CPU profile. Default false.
        The generated CPU profile is in <cpu.prof>. 
        Use "go tool pprof cpu.prof" to analyse the profile.

    -h : Optional. Display this help.

    -bs: Optional. Print BSR statistics.

    -v : Optional. Verbose: generate additional information files.
    
    -version : Optional. Display the version of this compiler`
	fmt.Println(msg)
}
