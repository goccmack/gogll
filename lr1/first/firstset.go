package first

import (
	"bytes"
	"fmt"
)

type FirstSet []string

func (this FirstSet) Add(sym ...string) (fs FirstSet, new bool) {
	fs = this
	for _, s := range sym {
		if !this.Contain(s) {
			fs = append(fs, s)
			new = true
		}
	}
	return
}

func (this FirstSet) Contain(sym string) bool {
	if this == nil {
		return false
	}
	for _, s := range this {
		if s == sym {
			return true
		}
	}
	return false
}

func (this FirstSet) Equal(that FirstSet) bool {
	if len(this) != len(that) {
		return false
	}
	for _, sym := range this {
		if !that.Contain(sym) {
			return false
		}
	}
	return true
}

/*
Return this - {sym}
*/
func (this FirstSet) Min(sym string) FirstSet {
	min := make(FirstSet, 0, len(this))
	for _, s := range this {
		if s != sym {
			min = append(min, s)
		}
	}
	return min
}

func (this FirstSet) Size() int {
	return len(this)
}

func (this FirstSet) String() string {
	w := new(bytes.Buffer)
	fmt.Fprintf(w, "{")
	for i, sym := range this {
		if i > 0 {
			fmt.Fprintf(w, ", %s", sym)
		} else {
			fmt.Fprintf(w, "%s", sym)
		}
	}
	fmt.Fprintf(w, "}")
	return w.String()
}
