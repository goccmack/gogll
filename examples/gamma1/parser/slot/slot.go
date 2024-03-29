
// Package slot is generated by gogll. Do not edit. 
package slot

import(
	"bytes"
	"fmt"
	
	"gamma1/parser/symbols"
)

type Label int

const(
	A0R0 Label = iota
	A0R1
	A1R0
	A1R1
	A2R0
	B0R0
	B0R1
	B1R0
	B1R1
	B1R2
	B2R0
	S0R0
	S0R1
	S0R2
	S0R3
	S1R0
	S1R1
	S1R2
	S1R3
)

type Slot struct {
	NT      symbols.NT
	Alt     int
	Pos     int
	Symbols symbols.Symbols
	Label 	Label
}

type Index struct {
	NT      symbols.NT
	Alt     int
	Pos     int
}

func GetAlternates(nt symbols.NT) []Label {
	alts, exist := alternates[nt]
	if !exist {
		panic(fmt.Sprintf("Invalid NT %s", nt))
	}
	return alts
}

func GetLabel(nt symbols.NT, alt, pos int) Label {
	l, exist := slotIndex[Index{nt,alt,pos}]
	if exist {
		return l
	}
	panic(fmt.Sprintf("Error: no slot label for NT=%s, alt=%d, pos=%d", nt, alt, pos))
}

func (l Label) EoR() bool {
	return l.Slot().EoR()
}

func (l Label) Head() symbols.NT {
	return l.Slot().NT
}

func (l Label) Index() Index {
	s := l.Slot()
	return Index{s.NT, s.Alt, s.Pos}
}

func (l Label) Alternate() int {
	return l.Slot().Alt
}

func (l Label) Pos() int {
	return l.Slot().Pos
}

func (l Label) Slot() *Slot {
	s, exist := slots[l]
	if !exist {
		panic(fmt.Sprintf("Invalid slot label %d", l))
	}
	return s
}

func (l Label) String() string {
	return l.Slot().String()
}

func (l Label) Symbols() symbols.Symbols {
	return l.Slot().Symbols
}

func (s *Slot) EoR() bool {
	return s.Pos >= len(s.Symbols)
}

func (s *Slot) String() string {
	buf := new(bytes.Buffer)
	fmt.Fprintf(buf, "%s : ", s.NT)
	for i, sym := range s.Symbols {
		if i == s.Pos {
			fmt.Fprintf(buf, "∙")
		}
		fmt.Fprintf(buf, "%s ", sym)
	}
	if s.Pos >= len(s.Symbols) {
		fmt.Fprintf(buf, "∙")
	}
	return buf.String()
}

var slots = map[Label]*Slot{ 
	A0R0: {
		symbols.NT_A, 0, 0, 
		symbols.Symbols{  
			symbols.T_0,
		}, 
		A0R0, 
	},
	A0R1: {
		symbols.NT_A, 0, 1, 
		symbols.Symbols{  
			symbols.T_0,
		}, 
		A0R1, 
	},
	A1R0: {
		symbols.NT_A, 1, 0, 
		symbols.Symbols{  
			symbols.T_2,
		}, 
		A1R0, 
	},
	A1R1: {
		symbols.NT_A, 1, 1, 
		symbols.Symbols{  
			symbols.T_2,
		}, 
		A1R1, 
	},
	A2R0: {
		symbols.NT_A, 2, 0, 
		symbols.Symbols{ 
		}, 
		A2R0, 
	},
	B0R0: {
		symbols.NT_B, 0, 0, 
		symbols.Symbols{  
			symbols.T_1,
		}, 
		B0R0, 
	},
	B0R1: {
		symbols.NT_B, 0, 1, 
		symbols.Symbols{  
			symbols.T_1,
		}, 
		B0R1, 
	},
	B1R0: {
		symbols.NT_B, 1, 0, 
		symbols.Symbols{  
			symbols.NT_B, 
			symbols.T_2,
		}, 
		B1R0, 
	},
	B1R1: {
		symbols.NT_B, 1, 1, 
		symbols.Symbols{  
			symbols.NT_B, 
			symbols.T_2,
		}, 
		B1R1, 
	},
	B1R2: {
		symbols.NT_B, 1, 2, 
		symbols.Symbols{  
			symbols.NT_B, 
			symbols.T_2,
		}, 
		B1R2, 
	},
	B2R0: {
		symbols.NT_B, 2, 0, 
		symbols.Symbols{ 
		}, 
		B2R0, 
	},
	S0R0: {
		symbols.NT_S, 0, 0, 
		symbols.Symbols{  
			symbols.T_0, 
			symbols.NT_A, 
			symbols.NT_B,
		}, 
		S0R0, 
	},
	S0R1: {
		symbols.NT_S, 0, 1, 
		symbols.Symbols{  
			symbols.T_0, 
			symbols.NT_A, 
			symbols.NT_B,
		}, 
		S0R1, 
	},
	S0R2: {
		symbols.NT_S, 0, 2, 
		symbols.Symbols{  
			symbols.T_0, 
			symbols.NT_A, 
			symbols.NT_B,
		}, 
		S0R2, 
	},
	S0R3: {
		symbols.NT_S, 0, 3, 
		symbols.Symbols{  
			symbols.T_0, 
			symbols.NT_A, 
			symbols.NT_B,
		}, 
		S0R3, 
	},
	S1R0: {
		symbols.NT_S, 1, 0, 
		symbols.Symbols{  
			symbols.T_0, 
			symbols.NT_A, 
			symbols.T_1,
		}, 
		S1R0, 
	},
	S1R1: {
		symbols.NT_S, 1, 1, 
		symbols.Symbols{  
			symbols.T_0, 
			symbols.NT_A, 
			symbols.T_1,
		}, 
		S1R1, 
	},
	S1R2: {
		symbols.NT_S, 1, 2, 
		symbols.Symbols{  
			symbols.T_0, 
			symbols.NT_A, 
			symbols.T_1,
		}, 
		S1R2, 
	},
	S1R3: {
		symbols.NT_S, 1, 3, 
		symbols.Symbols{  
			symbols.T_0, 
			symbols.NT_A, 
			symbols.T_1,
		}, 
		S1R3, 
	},
}

var slotIndex = map[Index]Label { 
	Index{ symbols.NT_A,0,0 }: A0R0,
	Index{ symbols.NT_A,0,1 }: A0R1,
	Index{ symbols.NT_A,1,0 }: A1R0,
	Index{ symbols.NT_A,1,1 }: A1R1,
	Index{ symbols.NT_A,2,0 }: A2R0,
	Index{ symbols.NT_B,0,0 }: B0R0,
	Index{ symbols.NT_B,0,1 }: B0R1,
	Index{ symbols.NT_B,1,0 }: B1R0,
	Index{ symbols.NT_B,1,1 }: B1R1,
	Index{ symbols.NT_B,1,2 }: B1R2,
	Index{ symbols.NT_B,2,0 }: B2R0,
	Index{ symbols.NT_S,0,0 }: S0R0,
	Index{ symbols.NT_S,0,1 }: S0R1,
	Index{ symbols.NT_S,0,2 }: S0R2,
	Index{ symbols.NT_S,0,3 }: S0R3,
	Index{ symbols.NT_S,1,0 }: S1R0,
	Index{ symbols.NT_S,1,1 }: S1R1,
	Index{ symbols.NT_S,1,2 }: S1R2,
	Index{ symbols.NT_S,1,3 }: S1R3,
}

var alternates = map[symbols.NT][]Label{ 
	symbols.NT_S:[]Label{ S0R0,S1R0 },
	symbols.NT_A:[]Label{ A0R0,A1R0,A2R0 },
	symbols.NT_B:[]Label{ B0R0,B1R0,B2R0 },
}

