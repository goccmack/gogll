package items

import (
	"bytes"
	"fmt"

	"github.com/goccmack/gogll/v3/ast"
	"github.com/goccmack/gogll/v3/lr1/basicprod"
)

type Items struct {
	List []*Item

	// key: item.haskHey
	idxMap map[string]int

	// key: prod id
	startItems map[string][]*Item
}

func NewItems(prods []*basicprod.Production) *Items {
	items := &Items{
		List:       make([]*Item, 0, 128),
		idxMap:     make(map[string]int),
		startItems: make(map[string][]*Item),
	}
	prodIdx := 0
	for _, prod := range prods {
		items.newItems(prod.Head, prod.Body, prodIdx)
		prodIdx++
	}
	return items
}

func (this *Items) StartItems(prodId string) []*Item {
	return this.startItems[prodId]
}

func (this *Items) newItems(head string, alt *ast.SyntaxAlternate, prodIdx int) {
	items := NewItem(head, alt, prodIdx)
	this.startItems[head] = append(this.startItems[head], items[0])
	for _, item := range items {
		if _, exist := this.idxMap[item.hashKey]; exist {
			panic(fmt.Sprintf("Duplicate item: %s", item))
		}
		this.List = append(this.List, item)
		this.idxMap[item.hashKey] = len(this.List) - 1
	}
}

func (this *Items) String() string {
	w := new(bytes.Buffer)
	for i, item := range this.List {
		fmt.Fprintf(w, "%d: %s\n", i, item)
		if item.NextItem != nil {
			fmt.Fprintf(w, "\tnext: %d\n\n", this.idxMap[item.NextItem.hashKey])
		} else {
			fmt.Fprintf(w, "\n")
		}
	}
	return w.String()
}
