package check

import (
	"fmt"
	"gogll/token"
	"os"
)

func fail(pos token.Pos, format string, a ...interface{}) {
	msg := fmt.Sprintf(format, a...)
	fmt.Printf("ERROR: %s at %s\n", msg, pos)
	os.Exit(1)
}
