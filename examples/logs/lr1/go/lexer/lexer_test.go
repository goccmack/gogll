package lexer

import (
	"fmt"
	"testing"
	"time"
	"unicode"
)

const REP = 1000_000_000

func _Test1(t *testing.T) {
	start := time.Now()
	for i := 0; i < REP; i++ {
		_ = nextState[0]('[')
	}
	fmt.Printf("Test1 %d ms\n", time.Now().Sub(start)/time.Millisecond)
}

func _Test2(t *testing.T) {
	vec := make([]int, 0, 2048)
	start := time.Now()
	for i := 0; i < 100_000; i++ {
		vec = append(vec, i)
	}
	fmt.Printf("Test2 %d mus\n", time.Now().Sub(start)/time.Microsecond)
}

func _Test3(t *testing.T) {
	start := time.Now()
	for i := 0; i < 100_000; i++ {
		unicode.IsLetter('9')
	}
	fmt.Printf("Test3 %d mus\n", time.Now().Sub(start)/time.Microsecond)
}

func _Test4(t *testing.T) {
	start := time.Now()
	for i := 0; i < 100_000_000; i++ {
		_ = nextState[1]('[')
	}
	fmt.Printf("Test1 %d ms\n", time.Now().Sub(start)/time.Millisecond)
}

func Test5(t *testing.T) {
	ptrn := []rune{'"'}
	start := time.Now()
	res := false
	for i := 0; i < 1000_000_000; i++ {
		for _, c := range ptrn {
			res = c == 'a'
		}
	}
	fmt.Printf("Test1 %d ms\n", time.Now().Sub(start)/time.Millisecond)
	fmt.Println(res)
}
