package parser

import (
	"bytes"
	"fmt"
	"gogll/ast"
	"gogll/frstflw"
	"gogll/gslot"
	"strings"
	"text/template"
)

func genTestSelect() string {
	tmpl, err := template.New("Test Select").Parse(testSelectTmpl)
	if err != nil {
		panic(err)
	}
	buf, data := new(bytes.Buffer), getTestSelectData()
	if err = tmpl.Execute(buf, data); err != nil {
		panic(err)
	}
	return buf.String()
}

type TestSelectData struct {
	TestSelect []*TSData
	Follow     []*TSData
}

type TSData struct {
	Label      string
	Conditions []*Condition
}

type Condition struct {
	Cond string
	Last bool
}

func getTestSelectData() *TestSelectData {
	return &TestSelectData{
		TestSelect: getTSData(),
		Follow:     getFollowData(),
	}
}

func getFollowData() (data []*TSData) {
	for _, nt := range ast.GetNonTerminals() {
		data = append(data, getFollowDataForNT(nt))
	}
	return
}

func getFollowDataForNT(nt string) *TSData {
	d := &TSData{
		Label:      nt,
		Conditions: getFollowConditions(nt),
	}
	// fmt.Printf("getFollowDataForNT(%s): %d\n", nt, len(d.Conditions))
	return d
}

func getTSData() (data []*TSData) {
	for _, s := range gslot.GetSlots() {
		data = append(data, getSlotTSData(s))
	}
	return
}

func getSlotTSData(s gslot.Label) *TSData {
	return &TSData{
		Label:      s.Label(),
		Conditions: getSlotTSConditions(s),
	}
}

func getSlotTSConditions(s gslot.Label) (data []*Condition) {
	ss := s.Symbols()[s.Pos:]
	frst := frstflw.FirstOfString(ss)
	for _, sym := range frst.Elements() {
		if sym != ast.Empty {
			data = append(data, &Condition{Cond: getSymbolCondition(sym)})
		}
	}
	if frst.Contain(ast.Empty) {
		data = append(data, getFollowConditions(s.Head)...)
	}
	data[len(data)-1].Last = true
	return
}

func getFollowConditions(nt string) (data []*Condition) {
	for _, sym := range frstflw.Follow(nt).Elements() {
		data = append(data, &Condition{Cond: getSymbolCondition(sym)})
	}
	data[len(data)-1].Last = true
	return
}

func getSymbolCondition(sym string) string {
	switch sym {
	case "any":
		return "true"
	case "letter":
		return "letter(r)"
	case "number":
		return "number(r)"
	case "upcase":
		return "upper(r)"
	case "lowcase":
		return "lower(r)"
	case "space":
		return "space(r)"
	}
	if strings.HasPrefix(sym, "not(") {
		set := strings.TrimSuffix(strings.TrimPrefix(sym, "not("), ")")
		return fmt.Sprintf(`not(r, %s)`, set)
	}
	return fmt.Sprintf("r == '%s'", sym)
}

const testSelectTmpl = `var testSelect = map[slot.Label]func()bool{ {{range $i, $ts := .TestSelect}}
	slot.{{$ts.Label}}:func()bool{
		return {{range $i, $c := $ts.Conditions}}{{$c.Cond}} {{if not $c.Last}}||{{end}}
		{{end}} },
{{end}} }

var follow = map[string]func()bool{ {{range $i, $flw := .Follow}}
	"{{$flw.Label}}":func()bool{
		return {{range $i, $c := $flw.Conditions}}{{$c.Cond}} {{if not $c.Last}}||{{end}}
	{{end}} },
{{end}} }`
