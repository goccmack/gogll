package golang

import (
	"bytes"
	"text/template"
)

func testFollowCode(nt string) string {
	conds := testFollowConditions(g.Follow(nt).Elements())
	// fmt.Printf("testFollow(%s) follow=s, conds = %s\n", nt, conds)
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

func testFollowConditions(follow []string) (conds []string) {
	for _, sym := range follow {
		conds = append(conds, getTestSelectCondition(sym))
	}
	return
}

const testFollowTemplate = `if {{range $i, $cond := .}}
			{{$cond}}{{end}}{
				stack.Pop(cU, cI, cN)
				L = L0
			}
`
