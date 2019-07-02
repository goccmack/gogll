package parser

import (
	"bytes"
	"fmt"
	"gogll/ast"
	"gogll/gslot"
	"text/template"
)

func codeAlt(rule *ast.Rule, altI int) string {
	symbols := rule.Alternates[altI].Symbols()
	if symbols[0] == ast.Empty {
		return codeTNEmpty(gslot.SlotLabel{rule.Head.Value(), altI, 1})
	}
	return codeNonEmpty(rule.Head.Value(), symbols, altI)
}

func codeNonEmpty(nt string, symbols []string, altI int) string {
	tmpl, err := template.New("codeString").Parse(nonEmpty)
	if err != nil {
		panic(err)
	}
	buf, data := new(bytes.Buffer), getNonEmptyData(nt, altI, symbols)
	if err := tmpl.Execute(buf, data); err != nil {
		panic(err)
	}
	return buf.String()
}

type DataNonEmpty struct {
	HeadCodeTN string
	Tail       []*TailSymbol
	TestFollow string
}

type TailSymbol struct {
	TestSelect string
	CodeTN     string
}

func getNonEmptyData(head string, altI int, symbols []string) *DataNonEmpty {
	data := &DataNonEmpty{
		HeadCodeTN: codeTN(symbols[0], gslot.SlotLabel{head, altI, 1}),
		TestFollow: testFollowConditions(head),
	}
	for i, sym := range symbols {
		if i == 0 {
			data.HeadCodeTN = codeTN(sym, gslot.SlotLabel{head, altI, 1})
		} else {
			symCode := getTailSymCode(head, sym, gslot.SlotLabel{head, altI, i + 1})
			data.Tail = append(data.Tail, symCode)
		}
	}
	return data
}

func getTailSymCode(nt, sym string, l gslot.SlotLabel) *TailSymbol {
	data := &TailSymbol{
		TestSelect: getTestSelectConditions(nt, sym),
		CodeTN:     codeTN(sym, l),
	}
	return data
}

// 	indent := "			"
// 	buf := new(bytes.Buffer)
// 	fmt.Fprintf(buf, codeTN(symbols[0], gslot.SlotLabel{rule.Head.Value(), altI, 1}))
// 	for i := 1; i < len(symbols); i++ {
// 		fmt.Fprintf(buf, "%s!(%s)\n", indent, getTestSelectForSymbol(rule.Head.Value(), symbols[i]))
// 		sl := gslot.SlotLabel{rule.Head.Value(), altI, i + 1}
// 		fmt.Fprintf(buf, "%s%s\n", indent, codeTN(symbols[i], sl))
// 	}
// 	fmt.Fprintf(buf, testFollowCode(rule.Head.Value()))
// 	return buf.String()
// }

const nonEmpty = `
			{{.HeadCodeTN}}
			{{range $i, $sym := .Tail}}if !({{$sym.TestSelect}}) {
				L = labels.L0
				break
			}
			{{$sym.CodeTN}}
			{{end}}
			if {{.TestFollow}} {
				pop(cU, cI, cN)
				L = labels.L0
				break
			}
`

func codeTNEmpty(sl gslot.SlotLabel) string {
	buf := new(bytes.Buffer)
	buf.WriteString("			cR = sppf.GetNodeE(cI)\n")
	fmt.Fprintf(buf, "			cN = sppf.GetNode(labels.%s, cN, cR)\n", sl.Label())
	buf.WriteString("			pop(cU, cI, cN)\n")
	buf.WriteString("			L = labels.L0\n")
	return buf.String()
}
