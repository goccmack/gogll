/*
Packge sppf implements a binarised shared packed parse forest
*/
package sppf

// Empty denotes the empty string. It is represented by the character, 'ϵ'
const Empty = "ϵ"

var (
	intermediateNodes = make(map[InKey]*IntermediateNode)
	packedNodes       = make(map[PNKey]*PackedNode)
	symbolNodes       = make(map[SNKey]*SymbolNode)
	// Key=slot label
	slots = map[int]*GrammarSlot{}
)

type GrammarSlot struct {
	Head   string
	EoR    bool
	FiR    bool
	String string
}

type SPPFNode interface {
	Equal(SPPFNode) bool
	GetExtent() (left, right int)
	GetLeftExtent() int
	GetRightExtent() int
	HasChild(L int, k int) bool
	SetChild(SPPFNode)
}

type Extent struct {
	Left  int
	Right int
}

var Dummy SPPFNode = nil

type Edge struct {
	from, to SPPFNode
}

/*
GetNode returns an SPPF node. L is the grammar slot label used by the parser. pos is the position
within slot L.
*/
func GetNode(L int, w, z SPPFNode) SPPFNode {
	if slots[L].FiR {
		return z
	}
	var y SPPFNode
	var j int
	k := z.GetLeftExtent()
	i := z.GetRightExtent()
	if w == Dummy {
		j = z.GetLeftExtent()
	} else {
		j = w.GetRightExtent()
	}
	if slots[L].EoR {
		y = getSymbolNode(slots[L].Head, j, i)
	} else {
		y = getIntermediateNode(L, j, i)
	}
	if !y.HasChild(L, k) {
		child := &PackedNode{
			PNKey: PNKey{
				SlotLabel:   L,
				RightExtent: k,
			},
			LeftChild:  w,
			RightChild: z,
		}
		y.SetChild(child)
		packedNodes[child.PNKey] = child
	}
	return y
}

func GetNodeE(i int) SPPFNode {
	return getSymbolNode(Empty, i, i)
}

func GetNodeT(symbol string, i, size int) *SymbolNode {
	return getSymbolNode(symbol, i, i+size)
}

func getIntermediateNode(L, j, i int) *IntermediateNode {
	sn := &IntermediateNode{
		InKey: InKey{
			SlotLabel: L,
			Extent:    Extent{j, i},
		},
	}
	if n, exist := intermediateNodes[sn.InKey]; !exist {
		intermediateNodes[sn.InKey] = sn
		return sn
	} else {
		return n
	}
}

func getSymbolNode(symbol string, i, j int) *SymbolNode {
	sn := &SymbolNode{
		SNKey: SNKey{
			Symbol: symbol,
			Extent: Extent{i, j},
		},
	}
	if n, exist := symbolNodes[sn.SNKey]; !exist {
		symbolNodes[sn.SNKey] = sn
		return sn
	} else {
		return n
	}
}

/*** IntermediateNode ***/

type IntermediateNode struct {
	InKey
	Child *PackedNode
}

type InKey struct {
	SlotLabel int
	Extent
}

func (sn *IntermediateNode) Equal(n SPPFNode) bool {
	n1, ok := n.(*IntermediateNode)
	if !ok {
		return false
	}
	return sn.SlotLabel == n1.SlotLabel && sn.Extent == n1.Extent
}

func (sn *IntermediateNode) GetExtent() (left, right int) {
	return sn.Extent.Left, sn.Extent.Right
}

func (sn *IntermediateNode) GetLeftExtent() int {
	return sn.Extent.Left
}

func (sn *IntermediateNode) GetRightExtent() int {
	return sn.Extent.Right
}

func (sn *IntermediateNode) HasChild(L, k int) bool {
	if sn.Child == nil {
		return false
	}
	return sn.Child.EqualPN(L, k)
}

func (sn *IntermediateNode) SetChild(pn SPPFNode) {
	sn.Child = pn.(*PackedNode)
}

/*** PackedNode ***/

type PackedNode struct {
	PNKey
	LeftChild  SPPFNode
	RightChild SPPFNode
}

type PNKey struct {
	SlotLabel   int
	RightExtent int
}

func (sn *PackedNode) Equal(n SPPFNode) bool {
	n1, ok := n.(*PackedNode)
	if !ok {
		return false
	}
	return sn.SlotLabel == n1.SlotLabel && sn.RightExtent == n1.RightExtent
}

func (sn *PackedNode) EqualPN(L, k int) bool {
	return sn.SlotLabel == L && sn.RightExtent == k
}

func (sn *PackedNode) GetExtent() (left, right int) {
	panic("not possible")
}

func (sn *PackedNode) GetLeftExtent() int {
	panic("not possible")
}

func (sn *PackedNode) GetRightExtent() int {
	return sn.RightExtent
}

func (sn *PackedNode) HasChild(L, k int) bool {
	panic("impossible")
}

func (sn *PackedNode) SetChild(pn SPPFNode) {
	panic("impossible")
}

/*** SymbolNode ***/

type SymbolNode struct {
	SNKey
	Child *PackedNode
}

type SNKey struct {
	Symbol string
	Extent
}

func (sn *SymbolNode) Equal(n SPPFNode) bool {
	n1, ok := n.(*SymbolNode)
	if !ok {
		return false
	}
	return sn.Symbol == n1.Symbol && sn.Extent == n1.Extent
}

func (sn *SymbolNode) GetExtent() (left, right int) {
	return sn.Extent.Left, sn.Extent.Right
}

func (sn *SymbolNode) GetLeftExtent() int {
	return sn.Extent.Left
}

func (sn *SymbolNode) GetRightExtent() int {
	return sn.Extent.Right
}

func (sn *SymbolNode) HasChild(L, k int) bool {
	if sn.Child == nil {
		return false
	}
	return sn.Child.EqualPN(L, k)
}

func (sn *SymbolNode) SetChild(pn SPPFNode) {
	sn.Child = pn.(*PackedNode)
}
