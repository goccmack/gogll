package utf8

import (
	"fmt"
	"unicode/utf8"
)

/*
DecodeRune returns the first rune in str as a string, together with the size of the first rune in bytes.
*/
func DecodeRune(str []byte) (string, int) {
	r, sz := utf8.DecodeRune(str)
	if r == utf8.RuneError {
		panic(fmt.Sprintf("Rune error: %s", str))
	}
	chr := RuneToString(r)
	return chr, sz
}

func DecodeRunes(str []byte) (ss []string) {
	for i := 0; i < len(str); {
		s, sz := DecodeRune(str[i:])
		i += sz
		if s == "\\" {
			s1, sz1 := DecodeRune(str[i:])
			i += sz1
			ss = append(ss, s1)
		} else {
			ss = append(ss, s)
		}
	}
	return
}

func RuneToString(r rune) string {
	buf := make([]byte, utf8.RuneLen(r))
	utf8.EncodeRune(buf, r)
	return string(buf)
}
