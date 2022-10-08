package symbols

import (
	"bytes"
	"fmt"
)

func (ss Symbols) String() string {
	w := new(bytes.Buffer)
	for _, s := range ss {
		fmt.Fprintf(w, "%s", s)
	}
	return w.String()
}
