package golang

import (
	"bytes"
	"fmt"
	"gogll/ast"
	"io/ioutil"
	"os"
	"path/filepath"
	"text/template"
)

type parserUtilsData struct {
	Package string
}

func genParserUtils(baseDir string) {
	data := &parserUtilsData{ast.GetPackage()}
	parserDir := filepath.Join(baseDir, "parser")
	genUtil(descriptors, parserDir, "descriptors", data)
	genUtil(gss, parserDir, "gss", data)
	genUtil(sppf, parserDir, "sppf", data)
}

func genUtil(src, baseDir, fname string, data *parserUtilsData) {
	tmpl, err := template.New(fname).Parse(src)
	if err != nil {
		failError(err)
	}
	buf := new(bytes.Buffer)
	if err = tmpl.Execute(buf, data); err != nil {
		failError(err)
	}
	if err := os.MkdirAll(baseDir, FilePerm); err != nil {
		failError(err)
	}
	file := filepath.Join(baseDir, fmt.Sprintf("%s.go", fname))
	if err := ioutil.WriteFile(file, buf.Bytes(), FilePerm); err != nil {
		failError(err)
	}
}

const descriptors = `
package parser

type Descriptors struct {
	// U, R Set
}

func NewDescriptors() *Descriptors {
	panic("implement me")
}

func (ds *Descriptors) Add(label int, u *GSSNode, i int, w SPPFNode) {
	panic("implement me")
}

func (ds *Descriptors) Empty() bool {
	panic("implement me")
}

func (ds *Descriptors) Remove() (L int, u *GSSNode, i int, w *SPPFNode) {
	panic("implement me")
}
`

const gss = `
package parser

type GSS struct{}

type GSSNode struct {
	// slot label
	L int
	// GSS node
	U	*GSSNode
	// Input position
	I     int
	// SPPF node
	W SPPFNode
}

func NewGSS() *GSS {
	panic("implement me")
}

func (gss *GSS) Create(slotLabel int, sn *GSSNode, iPos int, fn SPPFNode) *GSSNode {
	panic("implement me")
}

func (gss *GSS) Pop(u *GSSNode, i int, z SPPFNode){
	panic("implement me")
}
`

const sppf = `
package parser

type SPPF struct {}

func NewSPPF() *SPPF {
	panic("implement me")
}

type SPPFNode interface{}

var Dummy SPPFNode = nil

func (f *SPPF) Exist(symbol string, I int, m int) bool {
	panic("implement me")
}

func (f *SPPF) GetNode(L int, cN, cR SPPFNode) SPPFNode {
	panic("implement me")
}

func (f *SPPF) GetNodeE(I int) SPPFNode {
	panic("implement me")
}

func (f *SPPF) GetNodeT(symbol string, I int) SPPFNode {
	panic("implement me")
}

`
