package utils

// Escape returns a copy of s with all "\t\r\n\\ escaped
func Escape(s string) string {
	rs := []rune(s)
	ers := []rune{}
	for _, r := range rs {
		ers = append(ers, []rune(EscapeRune(r))...)
	}
	return string(ers)
}

// EscapeAll returns the elements of ss escaped
func EscapeAll(ss ...string) (ess []string) {
	for _, s := range ss {
		ess = append(ess, Escape(s))
	}
	return
}

// EscapeRune returns r escaped if in set: "\t\r\n\\
func EscapeRune(r rune) string {
	switch r {
	case '"':
		return "\\\""
	case '\\':
		return "\\\\"
	case '\r':
		return "\\r"
	case '\n':
		return "\\n"
	case '\t':
		return "\\t"
	}
	return string(r)
}

// StripEscape returns a copy of s with string literal escape sequences replaced
// by their escape characters.
func StripEscape(s string) (res string) {
	sr := []rune(s)
	runes := []rune{}
	for i := 0; i < len(sr); i++ {
		if sr[i] == '\\' {
			i++
			switch sr[i] {
			case 't':
				runes = append(runes, '\t')
			case 'r':
				runes = append(runes, '\r')
			case 'n':
				runes = append(runes, '\n')
			default:
				runes = append(runes, sr[i])
			}
		} else {
			runes = append(runes, sr[i])
		}
	}
	return string(runes)
}
