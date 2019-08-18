package slots

import (
	"bytes"
	"fmt"
	"gogll/cfg"
	"gogll/goutil/ioutil"
	"gogll/gslot"
	"os"
	"path/filepath"
)

func Gen() {
	buf := new(bytes.Buffer)
	for _, s := range gslot.GetSlots() {
		fmt.Fprintf(buf, "%s\n", s)
	}
	if err := ioutil.WriteFile(filepath.Join(cfg.BaseDir, "grammar_slots.txt"), buf.Bytes()); err != nil {
		fmt.Printf("Error writing grammar slots file: %s\n", err)
		os.Exit(1)
	}
}
