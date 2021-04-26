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

// Package gll generates Rust code for a GLL parser
package gll

import (
	"path/filepath"

	"github.com/goccmack/gogll/v3/ast"
	"github.com/goccmack/gogll/v3/frstflw"
	"github.com/goccmack/gogll/v3/gen/rust/gll/bsr"
	"github.com/goccmack/gogll/v3/gen/rust/gll/parser"
	"github.com/goccmack/gogll/v3/gen/rust/gll/slot"
	"github.com/goccmack/gogll/v3/gen/rust/gll/symbols"
	"github.com/goccmack/gogll/v3/gslot"
)

func Gen(parserDir string, g *ast.GoGLL, gs *gslot.GSlot, ff *frstflw.FF) {
	bsr.Gen(filepath.Join(parserDir, "bsr", "mod.rs"))
	symbols.Gen(filepath.Join(parserDir, "symbols", "mod.rs"), g)
	slot.Gen(filepath.Join(parserDir, "slot", "mod.rs"), g, gs, ff)
	parser.Gen(filepath.Join(parserDir, "mod.rs"), g, gs, ff)
}
