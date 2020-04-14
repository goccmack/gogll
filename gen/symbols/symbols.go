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

package symbols

import (
	"bytes"
	"fmt"
	"os"
	"path"

	"github.com/goccmack/gogll/ast"
	"github.com/goccmack/gogll/cfg"
	"github.com/goccmack/goutil/ioutil"
)

func Gen(g *ast.GoGLL) {
	if !cfg.Verbose {
		return
	}
	buf := new(bytes.Buffer)
	for _, sym := range g.GetSymbols() {
		fmt.Fprintf(buf, "%s\n", sym)
	}
	if err := ioutil.WriteFile(path.Join(cfg.BaseDir, "symbols.txt"), buf.Bytes()); err != nil {
		fail(err)
	}
}

func fail(err error) {
	fmt.Printf("Error writing symbols file: %s\n", err)
	os.Exit(1)
}
