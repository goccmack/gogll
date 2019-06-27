package golang

import (
	"bytes"
	"text/template"
)

func testFollowCode(nt string) string {
	conds := getTestSelectConditions(nt, 0, g.Follow(nt).Elements())
	tmpl, err := template.New("testFollowCode").Parse(testFollowTemplate)
	if err != nil {
		failError(err)
	}
	buf := new(bytes.Buffer)
	if err := tmpl.Execute(buf, conds); err != nil {
		failError(err)
	}
	return buf.String()
}

const testFollowTemplate = `if {{range $i, $cond := .}}
			{{$cond}}{{end}}{
				cU, cI, cN = pop()
				label = L0
			}
`
