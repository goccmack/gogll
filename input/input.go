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

/*
Package input wraps the gocc generated lexer to provide a slice of input tokens. Input[m]
*/
package input

import (
	"gogll2/lexer"
	"gogll2/token"
	"io/ioutil"
)

type Input struct {
	Buf      []byte
	All      []*token.Token
	Filtered []*token.Token
}

func New(srcFile string) *Input {
	buf, err := ioutil.ReadFile(srcFile)
	if err != nil {
		panic(err)
	}
	input := &Input{
		Buf: buf,
	}
	lex := lexer.NewLexer(buf)
	input.All = make([]*token.Token, 0, 1024)
	var tok *token.Token
	for tok := lex.Scan(); tok.Type != token.EOF; {
		input.All = append(input.All, tok)
	}
	input.All = append(input.All, tok)
	input.Filtered = filter(input.All)

	return input
}

func filter(all []*token.Token) []*token.Token {
	filtered := make([]*token.Token, 0, len(all))
	for _, tok := range all {
		if tok.Type != token.TokMap.Type("whitespace") &&
			tok.Type != token.TokMap.Type("comment") {
			filtered = append(filtered, tok)
		}
	}
	return filtered
}
