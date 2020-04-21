package runeset

import "sort"

// RuneSet holds a set of rune
type RuneSet struct {
	set map[rune]bool
}

// New returns a new RuneSet
func New() *RuneSet {
	return &RuneSet{
		set: make(map[rune]bool),
	}
}

// Add adds r to rs
func (rs *RuneSet) Add(r rune) {
	rs.set[r] = true
}

// Contains returns true iff r is in rs
func (rs *RuneSet) Contains(r rune) bool {
	return rs.set[r]
}

// Elements returns the runes in rs sorted in ascending order
func (rs *RuneSet) Elements() []rune {
	elements := make([]rune, 0, len(rs.set))
	for r := range rs.set {
		elements = append(elements, r)
	}
	sort.Slice(elements, func(i, j int) bool {
		return elements[i] < elements[j]
	})
	return elements
}

// Empty returns true iff rs has no elements
func (rs *RuneSet) Empty() bool {
	return len(rs.set) == 0
}

// Equal returns true iff rs and rs1 have the same elements
func (rs *RuneSet) Equal(rs1 *RuneSet) bool {
	if len(rs.set) != len(rs1.set) {
		return false
	}
	for r := range rs.set {
		if !rs1.set[r] {
			return false
		}
	}
	return true
}

// Intersection contains all elements common to rs and rs1
func (rs *RuneSet) Intersection(rs1 *RuneSet) *RuneSet {
	intersection := New()
	for r := range rs.set {
		if rs1.Contains(r) {
			intersection.Add(r)
		}
	}
	return intersection
}

// Subset returns true iff every element of rs is in rs1
func (rs *RuneSet) Subset(rs1 *RuneSet) bool {
	for r := range rs.set {
		if !rs1.set[r] {
			return false
		}
	}
	return true
}
