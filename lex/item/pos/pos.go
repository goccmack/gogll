// Package pos implements the position of a lexical dotted item
package pos

import (
	"bytes"
	"fmt"
)

// Pos contains the position of a lexical Item
type Pos struct {
	stack []int
}

// New returns a new position before the first symbol
func New() *Pos {
	p := &Pos{
		stack: make([]int, 0, 8),
	}
	return p.Push(0)
}

// Clone returns a deep copy of pos
func (pos *Pos) Clone() *Pos {
	clone := &Pos{
		stack: make([]int, len(pos.stack)),
	}
	copy(clone.stack[:len(pos.stack)], pos.stack[:len(pos.stack)])
	return clone
}

// Equal is true if the two stacks are the same
func (pos *Pos) Equal(other *Pos) bool {
	if len(pos.stack) != len(other.stack) {
		return false
	}
	for i := range pos.stack {
		if pos.stack[i] != other.stack[i] {
			return false
		}
	}
	return true
}

// Inc increments the value of pos.Top(). It returns pos for
// command chaining
func (pos *Pos) Inc() *Pos {
	pos.stack[len(pos.stack)-1]++
	return pos
}

// Len returns the number of items on the stack
func (pos *Pos) Len() int {
	return len(pos.stack)
}

// Peek returns the value of the item at height above the bottom of the stack.
// Pos remains unmodified
func (pos *Pos) Peek(height int) int {
	return pos.stack[height]
}

// Push pushes p onto the stack and returns a pointer to itself for command chainging.
func (pos *Pos) Push(p int) *Pos {
	pos.stack = append(pos.stack, p)
	return pos
}

// Pop removes the top n elements from the stack.
// Pop returns a self pointer for command chaining.
func (pos *Pos) Pop(n int) *Pos {
	pos.stack = pos.stack[:len(pos.stack)-n]
	return pos
}

func (pos *Pos) String() string {
	w := new(bytes.Buffer)
	fmt.Fprint(w, "[")
	for i, p := range pos.stack {
		if i > 0 {
			fmt.Fprint(w, ",")
		}
		fmt.Fprintf(w, "%d", p)
	}
	fmt.Fprint(w, "]")
	return w.String()
}

func (pos *Pos) Tail() *Pos {
	tail := pos.Clone()
	tail.stack = tail.stack[1:]
	return tail
}

// Top returns the value at the top of the stack without changing the stack
func (pos *Pos) Top() int {
	return pos.stack[len(pos.stack)-1]
}
