package parser

import (
	"bytes"
	"fmt"
	"gogll/frstflw"
	"gogll/gslot"
	"os"
	"strings"
	"text/template"
)

func (g *gen) genTestSelect() string {
	tmpl, err := template.New("Test Select").Parse(testSelectTmpl)
	if err != nil {
		panic(err)
	}
	buf, data := new(bytes.Buffer), g.getTestSelectData()
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

func (g *gen) getTestSelectData() *TestSelectData {
	return &TestSelectData{
		TestSelect: g.getTSData(),
		Follow:     g.getFollowData(),
	}
}

func (g *gen) getFollowData() (data []*TSData) {
	for _, nt := range g.g.GetNonTerminals() {
		data = append(data, g.getFollowDataForNT(nt))
	}
	return
}

func (g *gen) getFollowDataForNT(nt string) *TSData {
	d := &TSData{
		Label:      nt,
		Conditions: g.getFollowConditions(nt),
	}
	// fmt.Printf("getFollowDataForNT(%s): %d\n", nt, len(d.Conditions))
	return d
}

func (g *gen) getTSData() (data []*TSData) {
	for _, s := range g.gs.Slots() {
		data = append(data, g.getSlotTSData(s))
	}
	return
}

func (g *gen) getSlotTSData(s gslot.Label) *TSData {
	return &TSData{
		Label:      s.Label(),
		Conditions: g.getSlotTSConditions(s),
	}
}

func (g *gen) getSlotTSConditions(s gslot.Label) (data []*Condition) {
	// fmt.Printf("testselect.getSlotTSConditions(%s)\n", s)
	ss := s.Symbols()[s.Pos:]
	frst := g.ff.FirstOfString(ss)
	// fmt.Printf("  first: %s\n", frst)
	for _, sym := range frst.Elements() {
		if sym != frstflw.Empty {
			data = append(data, &Condition{Cond: getSymbolCondition(sym)})
		}
	}
	if frst.Contain(frstflw.Empty) {
		data = append(data, g.getFollowConditions(s.Head)...)
	}
	data[len(data)-1].Last = true
	return
}

func (g *gen) getFollowConditions(nt string) (data []*Condition) {
	// fmt.Printf("testselect.getFollowConditions(%s)\n", nt)
	flw := g.ff.Follow(nt)
	if flw.Len() == 0 {
		fmt.Printf("Production %s has empty follow set. It is never called\n", nt)
		os.Exit(1)
	}
	for _, sym := range flw.Elements() {
		data = append(data, &Condition{Cond: getSymbolCondition(sym)})
	}
	data[len(data)-1].Last = true
	return
}

func getSymbolCondition(sym string) string {
	// fmt.Printf("testselect.getSymbolCondition(%s)\n", sym)
	if sym == `'\\'` {
		panic(`'\\'`)
	}
	switch sym {
	case "any":
		return "true"
	case "letter":
		return "letter(r)"
	case "number":
		return "number(r)"
	case "upcase":
		return "upcase(r)"
	case "lowcase":
		return "lowcase(r)"
	case "space":
		return "space(r)"
	case "\\\"":
		return `r == '"'`
	}
	if strings.HasPrefix(sym, "not(") {
		set := strings.TrimSuffix(strings.TrimPrefix(sym, "not("), ")")
		return fmt.Sprintf(`not(r, "%s")`, set)
	}
	if strings.HasPrefix(sym, "anyof(") {
		set := strings.TrimSuffix(strings.TrimPrefix(sym, "anyof("), ")")
		return fmt.Sprintf(`anyof(r, "%s")`, set)
	}
	return fmt.Sprintf("r == '%s'", sym)
	// return fmt.Sprintf("r == '%s'", strings.TrimPrefix(sym, "\\"))
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
