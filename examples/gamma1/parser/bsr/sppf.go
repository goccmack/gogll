// [scott 2019] Elizabeth Scott, Adrian Johnstone, L. Thomas van Binsbergen
//              Derivation representation using binary subtree sets
//              Science of Computer Programming 175 (2019) 63–84

package bsr

import (
	"bytes"
	"fmt"
	"gamma1/parser/slot"
	"io"

	"github.com/goccmack/goutil/ioutil"
)

type SPPF struct {
	Root *NTNode
}

type Node interface {
	isNode()
}

type SymbolNode interface {
	isSymbolNode()
}

type IntermediateNode struct {
	Label      slot.Label
	Pivot      int
	LeftChild  *IntermediateNode
	RightChild *SymbolNode
}
type NTNode struct {
	Label       slot.Label
	Lext, Rext  int
	PackedNodes []*PackedNode
}

type PackedNode struct {
	Label      slot.Label
	Pivot      int
	LeftChild  *IntermediateNode
	RightChild *SymbolNode
}

func (*IntermediateNode) isNode() {}
func (*NTNode) isNode()           {}
func (*PackedNode) isNode()       {}

func (*NTNode) isSymbolNode() {}

type bldSPPF struct {
	spf          *SPPF
	extLeafNodes []Node
}

func (pf *Set) ToSPPF() *SPPF {
	rt := mkRootNode(pf)
	bldSPPF := &bldSPPF{
		spf: &SPPF{
			Root: rt,
		},
		extLeafNodes: []Node{rt},
	}
	for len(bldSPPF.extLeafNodes) > 0 {
		w := bldSPPF.extLeafNodes[len(bldSPPF.extLeafNodes)-1]
		bldSPPF.extLeafNodes = bldSPPF.extLeafNodes[:len(bldSPPF.extLeafNodes)-1]
		if nt, ok := w.(*NTNode); ok {
			bsts := pf.getNTSlot(nt.Label.Head(), nt.Lext, nt.Rext)
			for _, bst := range bsts {
				nt.PackedNodes = append(nt.PackedNodes, mkPN(bst))
			}
		}
	}
	return bldSPPF.spf
}

func mkPN(b BSR) *PackedNode {
	return &PackedNode{
		Label:      b.Label,
		Pivot:      b.pivot,
		LeftChild:  nil,
		RightChild: nil,
	}
}

func mkRootNode(pf *Set) *NTNode {
	rt := pf.GetRoots()[0]
	return &NTNode{
		Label: rt.Label,
		Lext:  rt.leftExtent,
		Rext:  rt.rightExtent,
	}
}

func (pn *PackedNode) ID() string {
	return fmt.Sprintf("\"%s:%s,%d\"", pn.Label.Head(), pn.Label.Symbols().String(), pn.Pivot)
}

func (nt *NTNode) ID() string {
	return fmt.Sprintf("\"%s,%d,%d\"", nt.Label.Head(), nt.Lext, nt.Rext)
}

//---- Dot ----

func (spf *SPPF) DotFile(file string) {
	w := new(bytes.Buffer)
	fmt.Fprintln(w, "digraph SPPF {")
	dotNT(w, spf.Root)
	fmt.Fprintln(w, "}")
	ioutil.WriteFile(file, w.Bytes())
}

func dotNT(w io.Writer, nt *NTNode) {
	fmt.Fprintln(w, nt.ID())
	for _, pn := range nt.PackedNodes {
		dotPN(w, pn)
		fmt.Fprintf(w, "%s -> %s\n", nt.ID(), pn.ID())
	}
	for i, pn := range nt.PackedNodes {
		if i > 0 {
			fmt.Fprint(w, ";")
		}
		fmt.Fprintf(w, "%s", pn.ID())
	}
	fmt.Fprintln(w)

}

func dotPN(w io.Writer, pn *PackedNode) {
	fmt.Fprintln(w, pn.ID())
}
