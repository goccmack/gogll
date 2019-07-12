package parser

import (
	"bytes"
	"fmt"
	"gogll/ast"
	"gogll/frstflw"
	"gogll/gslot"
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
	TS     []*TSData
	First  []*FFData
	Follow []*FFData
}

type FFData struct {
	Key   string
	Value string
}

type TSData struct {
	Label, NT string
}

func getTestSelectData() *TestSelectData {
	return &TestSelectData{
		TS:     getTSData(),
		First:  getFirstData(),
		Follow: getFollowData(),
	}
}

func getFirstData() (data []*FFData) {
	for _, s := range gslot.GetSlots() {
		data = append(data, getSlotFirstData(s))
	}
	return
}

func getFollowData() (data []*FFData) {
	for _, r := range ast.GetRules() {
		data = append(data, getNTFollowData(r.Head.Value()))
	}
	return
}

func getNTFollowData(nt string) *FFData {
	return &FFData{
		Key:   nt,
		Value: getSymbolsList(frstflw.Follow(nt).Elements()),
	}
}

func getSlotFirstData(s gslot.Label) *FFData {
	return &FFData{
		Key:   s.Label(),
		Value: getSymbolsList(frstflw.FirstOfString(s.Symbols()).Elements()),
	}
}

func getSymbolsList(symbols []string) string {
	buf := new(bytes.Buffer)
	for i, sym := range symbols {
		if i > 0 {
			buf.WriteString(",")
		}
		fmt.Fprintf(buf, "\"%s\"", sym)
	}
	return buf.String()
}

func getTSData() (data []*TSData) {
	for _, lbl := range gslot.GetSlots() {
		d := &TSData{
			Label: lbl.Label(),
			NT:    lbl.Head,
		}
		data = append(data, d)
	}
	return
}

// func getTSFirst(symbols []string) string {
// 	tmpl, err := template.New("TS first").Parse(firstTmpl)
// 	if err != nil {
// 		panic(err)
// 	}
// 	buf, data := new(bytes.Buffer), getFirstData()
// 	if err = tmpl.Execute(buf, data); err != nil {
// 		panic(err)
// 	}
// 	return buf.String()
// }

// func getFirstData() (data []*TestSelectData)

// func getTSFollow(nt string) string {
// 	println("implement me")
// }

const testSelectTmpl = `var testSelect = map[slot.Label]func(string)bool{ {{range $i, $ts := .TS}}
	slot.{{$ts.Label}}:func(x string)bool{
		return first[slot.{{$ts.Label}}].Contain(x) ||
			first[slot.{{$ts.Label}}].Contain(Empty) && follow["{{$ts.NT}}"].Contain(x)
	},{{end}}
}

var first = map[slot.Label]*stringset.StringSet {
	{{range $i, $f := .First}}slot.{{$f.Key}}:stringset.New({{$f.Value}}),
{{end}}}

var follow = map[string]*stringset.StringSet {
	{{range $i, $f := .Follow}}"{{$f.Key}}":stringset.New({{$f.Value}}),
{{end}}}
`

const firstTmpl = `
stringset.New({{range $i, $sym := .}}{{if i > 0}},{{end}}"{{$sym}}"{{end}})
`
