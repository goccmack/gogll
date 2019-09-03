//  Copyright 2019 Marius Ackerman
//
//  Licensed under the Apache License, Version 2.0 (the "License");
//  you may not use this file except in compliance with the License.
//  You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
//  Unless required by applicable law or agreed to in writing, software
//  distributed under the License is distributed on an "AS IS" BASIS,
//  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//  See the License for the specific language governing permissions and
//  limitations under the License.

package stringslice

type StringSlice []string

func (ss StringSlice) Add(e string) StringSlice {
	return append(ss, e)
}

func (ss StringSlice) Contain(e string) bool {
	for _, s := range ss {
		if s == e {
			return true
		}
	}
	return false
}

func (ss StringSlice) Equal(ss1 StringSlice) bool {
	if len(ss) != len(ss1) {
		return false
	}
	for i, s := range ss {
		if ss1[i] != s {
			return false
		}
	}
	return true
}

// Find returns a list of indices of ss which contain s.
// Find returns an empty slice if ss does not contain s.
func (ss StringSlice) Find(s string) (indices []int) {
	for i, s1 := range ss {
		if s1 == s {
			indices = append(indices, i)
		}
	}
	return
}
