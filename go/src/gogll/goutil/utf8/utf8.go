package utf8

import (
	"fmt"
	"unicode/utf8"
)

/*
DecodeRune returns the first rune in str as a string, together with the size of the first rune in bytes.
*/
// func DecodeRune(str []byte) (string, int) {
// 	r, sz := utf8.DecodeRune(str)
// 	if r == utf8.RuneError {
// 		panic(fmt.Sprintf("Rune error: %r", str))
// 	}
// 	chr := RuneToString(r)
// 	return chr, sz
// }

func DecodeRunes(str string) (rs []rune, err error) {
	for i := 0; i < len(str); {
		r, sz := utf8.DecodeRune([]byte(str[i:]))
		if r == utf8.RuneError {
			return nil, fmt.Errorf("Rune error in %s", str[i:])
		}
		rs = append(rs, r)
		i += sz
	}
	return
}
