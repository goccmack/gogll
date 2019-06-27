package slots

import (
	"bytes"
	"fmt"
	"gogll/gslot"
	"io/ioutil"
	"os"
)

func Gen(dir string, perm os.FileMode) {
	os.MkdirAll(dir, perm)
	buf := new(bytes.Buffer)
	for _, s := range gslot.GetSlots() {
		fmt.Fprintf(buf, "%s\n", s)
	}
	if err := ioutil.WriteFile(dir+"/grammar_slots.txt", buf.Bytes(), perm); err != nil {
		fmt.Printf("Error writing grammar slots file: %s\n", err)
		os.Exit(1)
	}
}
