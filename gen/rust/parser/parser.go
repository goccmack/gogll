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

// Package parser generates Rust code for a GLL parser
package parser

import (
	"path/filepath"

	"github.com/goccmack/gogll/ast"
	"github.com/goccmack/gogll/frstflw"
	"github.com/goccmack/gogll/gen/rust/parser/parser"
	"github.com/goccmack/gogll/gen/rust/parser/slot"
	"github.com/goccmack/gogll/gen/rust/parser/symbols"
	"github.com/goccmack/gogll/gslot"
	"github.com/goccmack/gogll/im/tokens"
)

func Gen(parserDir string, g *ast.GoGLL, gs *gslot.GSlot, ff *frstflw.FF, ts *tokens.Tokens) {
	// bsr.Gen(filepath.Join(parserDir, "bsr", "mod.rs"))
	symbols.Gen(filepath.Join(parserDir, "symbols", "mod.rs"), g)
	slot.Gen(filepath.Join(parserDir, "slot", "mod.rs"), g, gs, ff)
	parser.Gen(filepath.Join(parserDir, "mod.rs"), g, gs, ff, ts)
}
