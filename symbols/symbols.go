package symbols

type Symbols []string

func (ss Symbols) Empty() bool {
	return len(ss) == 0
}
