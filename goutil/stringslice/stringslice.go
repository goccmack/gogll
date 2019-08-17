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
