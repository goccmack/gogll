package descriptors

type Descriptor struct {
}

type Set = []Descriptor

func Empty() Set {
	return Set{}
}
