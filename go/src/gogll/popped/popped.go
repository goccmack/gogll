package popped

import (
	"gogll2/gss"
	"gogll2/sppf"
)

type Element struct {
	GSSNode  gss.Node
	SPPFNode sppf.Node
}

type Set []*Element

func Empty() Set {
	return Set{}
}
