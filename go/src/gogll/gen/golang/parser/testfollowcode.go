package parser

import "bytes"

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

