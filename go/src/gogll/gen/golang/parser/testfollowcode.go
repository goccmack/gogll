package parser

import "bytes"

// func testFollowCode(nt string) string {
// 	conds := testFollowConditions(g.Follow(nt).Elements())
// 	// fmt.Printf("testFollow(%s) follow=s, conds = %s\n", nt, conds)
// 	tmpl, err := template.New("testFollowCode").Parse(testFollowTemplate)
// 	if err != nil {
// 		failError(err)
// 	}
// 	buf := new(bytes.Buffer)
// 	if err := tmpl.Execute(buf, conds); err != nil {
// 		failError(err)
// 	}
// 	return buf.String()
// }

// func testFollowConditions(follow []string) (conds []string) {
// 	for _, sym := range follow {
// 		conds = append(conds, getTestSelectCondition(sym))
// 	}
// 	return
// }

func testFollowConditions(nt string) string {
	follow := g.Follow(nt).Elements()
	buf := new(bytes.Buffer)
	for i, sym := range follow {
		buf.WriteString(getTestSelectCondition(sym))
		if i < len(follow)-1 {
			buf.WriteString(" || \n")
		}
	}
	return buf.String()
}

// const testFollowTemplate = `if {{range $i, $cond := .}}
// 			{{$cond}}{{end}}{
// 				pop(cU, cI, cN)
// 				L = labels.L0
// 			}
// `
