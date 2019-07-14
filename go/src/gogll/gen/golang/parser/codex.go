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
		Slots:      getSlotsData(nt, alt, altI),
	}
	return d
}

func getSlotsData(nt string, alt *ast.Alternate, altI int) (data []*SlotData) {
	if alt.Empty() {
		data = append(data, &SlotData{Empty: true})
		return
	}
	for i, sym := range alt.Symbols() {
		data = append(data, getSlotData(nt, altI, sym, i))
	}
	return
}

func getSlotData(nt string, altI int, symbol string, pos int) *SlotData {
	// fmt.Printf("getSlotCode: %s\n", gslot.Label{nt, altI, pos}.String())
	preLabel, postLabel := gslot.Label{nt, altI, pos}, gslot.Label{nt, altI, pos + 1}
	sd := &SlotData{
		PreLabel:  preLabel.Label(),
		PostLabel: postLabel.Label(),
		Comment:   postLabel.String(),
		Head:      nt,
	}
	sd.IsNT = !ast.IsTerminal(symbol)
	return sd
}

type SlotData struct {
	PreLabel  string
	PostLabel string
	Comment   string
	Empty     bool
	IsNT      bool
	Head      string
}

const altCodeTmpl = `		case slot.{{.AltLabel}}: // {{.AltComment}}{{range $i, $slot := .Slots}}
			{{if $i}}if !testSelect[slot.{{$slot.PreLabel}}](nextI){ 
				break 
			}
			{{end}}{{if $slot.Empty}}bsr.AddEmpty({{$slot.PostLabel}},cI, cI, cI)
			{{else if $slot.IsNT}}call(slot.{{.PostLabel}}, cU, cI)
case slot.{{$slot.PostLabel}}: // {{$slot.Comment}} 
			{{else}}bsr.Add(slot.{{$slot.PostLabel}}, cU, cI, cI+sz)
			cI += sz 
			nextI, _, sz = decodeRune(I[cI:]){{end}}{{end}}
			if follow["{{.NT}}"].Contain(nextI){
				rtn("{{.NT}}", cU, cI)
			}
	`
