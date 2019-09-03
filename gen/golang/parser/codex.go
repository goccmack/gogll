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
	"gogll/ast"
	"gogll/gslot"
	"text/template"
)

func genAlternatesCode() string {
	buf := new(bytes.Buffer)
	for _, nt := range ast.GetNonTerminals() {
		rule := ast.GetRule(nt)
		for i, alt := range rule.Alternates {
			buf.WriteString(getAlternateCode(nt, alt, i))
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
}

func getAlternateCode(nt string, alt *ast.Alternate, altI int) string {
	tmpl, err := template.New("Alternate").Parse(altCodeTmpl)
	if err != nil {
		panic(err)
	}
	buf, data := new(bytes.Buffer), getAltData(nt, alt, altI)
	if err = tmpl.Execute(buf, data); err != nil {
		panic(err)
	}
	return buf.String()
}

func getAltData(nt string, alt *ast.Alternate, altI int) *AltData {
	L := gslot.Label{nt, altI, 0}
	d := &AltData{
		NT:         nt,
		AltLabel:   L.Label(),
		AltComment: L.String(),
		Empty:      alt.Empty(),
	}
	if !alt.Empty() {
		d.Slots = getSlotsData(nt, alt, altI)
	}
	return d
}

func getSlotsData(nt string, alt *ast.Alternate, altI int) (data []*SlotData) {
	for i, sym := range alt.Symbols() {
		// fmt.Printf("getSlotsData(%s) %s\n", nt, getSlotData(nt, altI, sym, i))
		data = append(data, getSlotData(nt, altI, sym, i))
	}
	return
}

func getSlotData(nt string, altI int, symbol string, pos int) *SlotData {
	preLabel, postLabel := gslot.Label{nt, altI, pos}, gslot.Label{nt, altI, pos + 1}
	sd := &SlotData{
		AltLabel:  gslot.Label{nt, altI, 0}.Label(),
		PreLabel:  preLabel.Label(),
		PostLabel: postLabel.Label(),
		Comment:   postLabel.String(),
		Head:      nt,
	}
	sd.IsNT = !ast.IsTerminal(symbol)
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
			bsr.AddEmpty(slot.{{.AltLabel}},cI)
		{{else}}{{range $i, $slot := .Slots}}
			{{if $i}}if !testSelect[slot.{{$slot.PreLabel}}](){ 
				parseError(slot.{{$slot.PreLabel}}, cI)
				break 
			}
			{{end}}
			{{if $slot.IsNT}}call(slot.{{$slot.PostLabel}}, cU, cI)
case slot.{{$slot.PostLabel}}: // {{$slot.Comment}} 
			{{else}}bsr.Add(slot.{{$slot.PostLabel}}, cU, cI, cI+sz)
			cI += sz 
			nextI, r, sz = decodeRune(I[cI:]){{end}}{{end}}{{end}}
			if follow["{{.NT}}"](){
				rtn("{{.NT}}", cU, cI)
			}
	`
