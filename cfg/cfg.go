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

// Package cfg reads the commandline options
package cfg

import (
	"flag"
	"fmt"
	"os"
	"path"
)

// Version is the version of this compiler
const Version = "v3.2.2"

var (
	BaseDir string
	SrcFile string
	Verbose bool

	All        = flag.Bool("a", false, "Regenerate all files")
	BSRStats   = flag.Bool("bs", false, "Print BSR stats")
	help       = flag.Bool("h", false, "Print help")
	CPUProfile = flag.Bool("CPUProf", false, "Generate CPU profile")
	outDir     = flag.String("o", "", "")
	verbose    = flag.Bool("v", false, "Verbose")
	version    = flag.Bool("version", false, "Version")

	Go   = flag.Bool("go", true, "Generate Go code")
	Rust = flag.Bool("rust", false, "Generate Rust code")

	target = flag.String("t", "go", "Target Language")

	GLL               = flag.Bool("gll", true, "Generate GLL parser")
	Knuth             = flag.Bool("knuth", false, "Generate Knuth LR(1) parser")
	Pager             = flag.Bool("pager", false, "Generate Pager's PGM parser")
	AutoResolveLRConf = flag.Bool("resolve_conflicts", false, "Auto resolve LR(1) conflicts")
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
	getParserType()
	if *Rust {
		fmt.Printf("Version %s does not support Rust\n", Version)
		fmt.Println("Please log an issue if you need Rust support")

		*Go = false
	}
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

func getParserType() {
	if *Pager || *Knuth {
		*GLL = false
	}
	if *Pager && *Knuth {
		fail("Only one of pager or knuth may be selected")
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
	msg :=
		`use: gogll -h
    for help, or

use: gogll -version
    to display the version of goggl, or

use: gogll [-a][-v] [-CPUProf] [-o <out dir>] [-go] [-rust] [-gll] [-pager] [-knuth] [-resolve_conflicts] <source file>
    to generate a lexer and parser.

    <source file>: Mandatory. Name of the source file to be processed. 
        If the file extension is ".md" the bnf is extracted from markdown code 
        segments enclosed in triple backticks.
    
    -a: Optional. Regenerate all files.
        WARNING: This may destroy user editing in the LR(1) AST.
        Default: false
         
    -v: Optional. Produce verbose output, including first and follow sets,
        LR(1) sets and lexer FSA sets.
    
    -o <out dir>: Optional. The directory to which code will be generated.
                  Default: the same directory as <source file>.
                  
    -go: Optional. Generate Go code.
          Default: true, but false if -rust is selected

    -rust: Optional. Generate Rust code.
           Default: false
           
    -gll: Optional. Generate a GLL parser.
          Default true. False if -knuth or -pager is selected.
                  
    -knuth: Optional. Generate a Knuth LR(1) parser
            Default false

    -pager: Optional. Generate a Pager PGM LR(1) parser.
            Default false

    -resolve_conflicts: Optional. Automatically resolve LR(1) conflicts.
            Default: false. Only used when generating LR(1) parsers.
    
    -bs: Optional. Print BSR statistics (GLL only).
    
    -CPUProf : Optional. Generate a CPU profile. Default false.
        The generated CPU profile is in <cpu.prof>. 
        Use "go tool pprof cpu.prof" to analyse the profile.`

	fmt.Println(msg)
}
