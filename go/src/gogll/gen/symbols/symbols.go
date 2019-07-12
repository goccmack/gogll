package symbols

import (
	"bytes"
	"fmt"
	"gogll/ast"
	"gogll/cfg"
	"gogll/goutil/ioutil"
	"os"
	"path"
)

func Gen() {
	buf := new(bytes.Buffer)
	for _, sym := range ast.GetSymbols() {
		fmt.Fprintf(buf, "%s\n", sym)
	}
	if err := ioutil.WriteFile(path.Join(cfg.BaseDir, "symbols.txt"), buf.Bytes()); err != nil {
		fail(err)
	}
}

func fail(err error) {
	fmt.Printf("Error writing symbols file: %s\n", err)
	os.Exit(1)
}
