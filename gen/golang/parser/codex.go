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

package parser

import (
	"bytes"
	"text/template"

	"github.com/goccmack/gogll/ast"
	"github.com/goccmack/gogll/gslot"
)

func (g *gen) genAlternatesCode() string {
	buf := new(bytes.Buffer)
	for _, nt := range g.g.NonTerminals.ElementsSorted() {
		rule := g.g.GetSyntaxRule(nt)
		for i, alt := range rule.Alternates {
			buf.WriteString(g.getAlternateCode(nt, alt, i))
		}
	}
	return buf.String()
}

type AltData struct {
	NT         string
	AltLabel   string
	AltComment string
	Empty      bool
	Slots      []*SlotData
	LastSlot   *SlotData
}

func (g *gen) getAlternateCode(nt string, alt *ast.SyntaxAlternate, altI int) string {
	// fmt.Printf("codes.genAlternateCode: %s altI %d\n", nt, altI)
	tmpl, err := template.New("Alternate").Parse(altCodeTmpl)
	if err != nil {
		panic(err)
	}
	buf, data := new(bytes.Buffer), g.getAltData(nt, alt, altI)
	if err = tmpl.Execute(buf, data); err != nil {
		panic(err)
	}
	return buf.String()
}

func (g *gen) getAltData(nt string, alt *ast.SyntaxAlternate, altI int) *AltData {
	// fmt.Printf("codex.getAltData %s[%d]\n", nt, altI)
	L := gslot.NewLabel(nt, altI, 0, g.gs, g.ff)
	d := &AltData{
		NT:         nt,
		AltLabel:   L.Label(),
		AltComment: L.String(),
		Empty:      alt.Empty(),
	}
	if !alt.Empty() {
		d.Slots = g.getSlotsData(nt, alt, altI)
		d.LastSlot = d.Slots[len(d.Slots)-1]
	}
	return d
}

func (g *gen) getSlotsData(nt string, alt *ast.SyntaxAlternate, altI int) (data []*SlotData) {
	for i, sym := range alt.Symbols {
		// fmt.Printf("getSlotsData(%s) %s\n", nt, getSlotData(nt, altI, sym, i))
		data = append(data, g.getSlotData(nt, altI, sym.String(), i))
	}
	return
}

func (g *gen) getSlotData(nt string, altI int, symbol string, pos int) *SlotData {
	preLabel := gslot.NewLabel(nt, altI, pos, g.gs, g.ff)
	postLabel := gslot.NewLabel(nt, altI, pos+1, g.gs, g.ff)
	sd := &SlotData{
		AltLabel:  gslot.NewLabel(nt, altI, 0, g.gs, g.ff).Label(),
		PreLabel:  preLabel.Label(),
		PostLabel: postLabel.Label(),
		Comment:   postLabel.String(),
		Head:      nt,
	}
	sd.IsNT = !g.g.Terminals.Contain(symbol)
	// fmt.Printf("getSlotData: altlabel:%s, pre:%s, post:%s\n",
	// 	sd.AltLabel, sd.PreLabel, sd.PostLabel)
	return sd
}

type SlotData struct {
	AltLabel  string
	PreLabel  string
	PostLabel string
	Comment   string
	IsNT      bool
	Head      string
}

const altCodeTmpl = `		case slot.{{.AltLabel}}: // {{.AltComment}}{{if .Empty}}
			p.bsrSet.AddEmpty(slot.{{.AltLabel}},cI)
		{{else}}{{range $i, $slot := .Slots}}
			{{if $i}}if !p.testSelect(slot.{{$slot.PreLabel}}){ 
				p.parseError(slot.{{$slot.PreLabel}}, p.cI, first[slot.{{$slot.PreLabel}}])
				break 
			}
			{{end}}
			{{if $slot.IsNT}}p.call(slot.{{$slot.PostLabel}}, cU, p.cI)
case slot.{{$slot.PostLabel}}: // {{$slot.Comment}} 
			{{else}}p.bsrSet.Add(slot.{{$slot.PostLabel}}, cU, p.cI, p.cI+1)
			p.cI++ {{end}}{{end}}{{end}}
			if p.follow(symbols.NT_{{.NT}}) {
				p.rtn(symbols.NT_{{.NT}}, cU, p.cI)
			} else { 
				p.parseError(slot.{{.AltLabel}}, p.cI, followSets[symbols.NT_{{.NT}}])
			}
	`
