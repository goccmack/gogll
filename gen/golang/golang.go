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

// Package golang controls the Go code generation
package golang

import (
	"path"

	"github.com/goccmack/gogll/ast"
	"github.com/goccmack/gogll/cfg"
	"github.com/goccmack/gogll/frstflw"
	"github.com/goccmack/gogll/gen/golang/lexer"
	"github.com/goccmack/gogll/gen/golang/parser"
	"github.com/goccmack/gogll/gen/golang/token"
	"github.com/goccmack/gogll/gslot"
	"github.com/goccmack/gogll/im/tokens"
	lexitems "github.com/goccmack/gogll/lex/items"
)

func Gen(g *ast.GoGLL, gs *gslot.GSlot, ff *frstflw.FF, ls *lexitems.Sets, ts *tokens.Tokens) {
	token.Gen(g, ts)
	lexer.Gen(path.Join(cfg.BaseDir, "lexer"), g, ls, ts)
	parser.Gen(path.Join(cfg.BaseDir, "parser"), g, gs, ff, ts)
}
