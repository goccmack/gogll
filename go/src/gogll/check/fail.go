package check

import (
	"fmt"
	"os"
)

func fail(msg string) {
	fmt.Println("ERROR: ", msg)
	os.Exit(1)
}
