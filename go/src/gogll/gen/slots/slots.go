package slots

import (
	"bytes"
	"fmt"
	"gogll/goutil/ioutil"
	"gogll/gslot"
	"os"
)

func Gen(dir string) {
	buf := new(bytes.Buffer)
	for _, s := range gslot.GetSlots() {
		fmt.Fprintf(buf, "%s\n", s)
	}
	if err := ioutil.WriteFile(dir+"/grammar_slots.txt", buf.Bytes()); err != nil {
		fmt.Printf("Error writing grammar slots file: %s\n", err)
		os.Exit(1)
	}
}
