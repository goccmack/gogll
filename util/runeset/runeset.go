/*
Copyright 2020 Marius Ackerman

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

  http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Package runeset implements a set of runes
package runeset

import (
	"bytes"
	"fmt"
	"sort"
)

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

func (rs *RuneSet) String() string {
	w := new(bytes.Buffer)
	fmt.Fprintf(w, "[]rune{")
	for i, r := range rs.Elements() {
		if i > 0 {
			fmt.Fprint(w, ",")
		}
		fmt.Fprintf(w, "'%s'", escape(r))
	}
	fmt.Fprint(w, "}")
	return w.String()
}

func escape(r rune) string {
	switch r {
	case '\\':
		return "\\\\"
	case '\'':
		return "\\'"
	case '\r':
		return "\\r"
	case '\n':
		return "\\n"
	case '\t':
		return "\\t"
	}
	return string(r)
}
