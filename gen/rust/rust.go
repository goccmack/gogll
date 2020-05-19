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

// Package rust generates code for the Rust language
package rust

import (
	"path"
	"path/filepath"

	"github.com/goccmack/gogll/ast"
	"github.com/goccmack/gogll/cfg"
	"github.com/goccmack/gogll/frstflw"
	"github.com/goccmack/gogll/gen/rust/lexer"
	"github.com/goccmack/gogll/gen/rust/parser"
	"github.com/goccmack/gogll/gen/rust/token"
	"github.com/goccmack/gogll/gslot"
	"github.com/goccmack/gogll/im/tokens"
	"github.com/goccmack/gogll/lex/items"
)

func Gen(g *ast.GoGLL, gs *gslot.GSlot, ff *frstflw.FF, ls *items.Sets, ts *tokens.Tokens) {
	token.Gen(filepath.Join(cfg.BaseDir, "src", "token", "mod.rs"), ts)
	lexer.Gen(path.Join(cfg.BaseDir, "src", "lexer", "mod.rs"), g, ls, ts)
	if len(g.SyntaxRules) > 0 {
		parser.Gen(path.Join(cfg.BaseDir, "src", "parser"), g, gs, ff, ts)
	}
}
