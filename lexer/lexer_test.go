package lexer

import (
	"testing"
)

// TestLoadMd_WithLanguageTag tests that language tags after opening backticks
// are properly ignored (replaced with spaces) in markdown code blocks.
func TestLoadMd_WithLanguageTag(t *testing.T) {
	tests := []struct {
		name string
		src  string
		exp  string
	}{
		{
			name: "code block with go language tag",
			src:  "```go\nvar code int\n```",
			exp:  "     \nvar code int\n   ",
		},
		{
			name: "code block with gogll language tag",
			src:  "```gogll\nsymbol\n  : not \"\\n\"\n```",
			exp:  "        \nsymbol\n  : not \"\\n\"\n   ",
		},
		{
			name: "multiple code blocks with language tags",
			src:  "text\n```gogll\nlexicon\n```\ntext\n```gogll\nsyntax\n```",
			exp:  "    \n        \nlexicon\n   \n    \n        \nsyntax\n   ",
		},
		{
			name: "code block without language tag",
			src:  "```\ncode\n```",
			exp:  "   \ncode\n   ",
		},
		{
			name: "text outside code blocks becomes spaces",
			src:  "some text\n```\ncode\n```\nmore text",
			exp:  "         \n   \ncode\n   \n         ",
		},
		{
			name: "language tag with info string",
			src:  "```gogll lexical rules\ncode\n```",
			exp:  "                      \ncode\n   ",
		},
		{
			name: "empty language tag line",
			src:  "```   \ncode\n```",
			exp:  "      \ncode\n   ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := []rune(tt.src)
			loadMd(got)

			if string(got) != tt.exp {
				t.Errorf("loadMd() mismatch\nsrc: %q\nexp: %q\ngot: %q",
					tt.src, tt.exp, string(got))

				// Show character-by-character comparison for debugging
				t.Logf("\nCharacter comparison:")
				maxLen := len(got)
				if len(tt.exp) > maxLen {
					maxLen = len(tt.exp)
				}
				for i := 0; i < maxLen; i++ {
					expChar := ' '
					gotChar := ' '
					if i < len(tt.exp) {
						expChar = rune(tt.exp[i])
					}
					if i < len(got) {
						gotChar = rune(got[i])
					}
					if expChar != gotChar {
						t.Logf("  [%d] exp %q (%d), got %q (%d)",
							i, expChar, expChar, gotChar, gotChar)
					}
				}
			}
		})
	}
}

// TestLoadMd_PreservesCodeContent tests that code inside code blocks
// is preserved exactly as-is, while text outside is replaced with spaces.
func TestLoadMd_PreservesCodeContent(t *testing.T) {
	src := "# Heading\nSome text before.\n```go\nfunc main() {\n\tfmt.Println(\"hello\")\n}\n```\nSome text after."
	exp := "         \n                 \n     \nfunc main() {\n\tfmt.Println(\"hello\")\n}\n   \n                "

	got := []rune(src)
	loadMd(got)

	if string(got) != exp {
		t.Errorf("loadMd() did not preserve code content correctly\nexp: %q\ngot: %q", exp, string(got))
	}
}

// TestLoadMd_PreservesNewlines tests that newlines are preserved throughout
// the markdown file, both inside and outside code blocks.
func TestLoadMd_PreservesNewlines(t *testing.T) {
	src := "line1\nline2\n```\ncode1\ncode2\n```\nline3\nline4"
	got := []rune(src)
	loadMd(got)

	// Count newlines in src and output
	srcNewlines := 0
	for _, r := range "line1\nline2\n```\ncode1\ncode2\n```\nline3\nline4" {
		if r == '\n' {
			srcNewlines++
		}
	}

	gotNewlines := 0
	for _, r := range got {
		if r == '\n' {
			gotNewlines++
		}
	}

	if srcNewlines != gotNewlines {
		t.Errorf("loadMd() changed number of newlines:\nexp %4d: %q\ngot %4d: %q",
			srcNewlines, "line1\nline2\n```\ncode1\ncode2\n```\nline3\nline4", gotNewlines, got)
	}
}

// TestLoadMd_EdgeCases tests edge cases and boundary conditions.
func TestLoadMd_EdgeCases(t *testing.T) {
	tests := []struct {
		name string
		src  string
		exp  string
	}{
		{
			name: "empty src",
			src:  "",
			exp:  "",
		},
		{
			name: "only backticks",
			src:  "```",
			exp:  "   ",
		},
		{
			name: "incomplete code block",
			src:  "```go\ncode",
			exp:  "     \ncode",
		},
		{
			name: "nested-looking backticks in code",
			src:  "```\ncode with ``` backticks\n```",
			exp:  "   \ncode with              \n   ",
		},
		{
			name: "language tag without newline",
			src:  "```go",
			exp:  "     ",
		},
		{
			name: "unicode in language tag",
			src:  "```gö\ncode\n```",
			exp:  "     \ncode\n   ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			src := tt.src
			got := []rune(src)
			loadMd(got)

			if string(got) != tt.exp {
				t.Errorf("loadMd() edge case mismatch\nsrc: %q\nexp: %q\ngot: %q",
					tt.src, tt.exp, string(got))
			}
		})
	}
}
