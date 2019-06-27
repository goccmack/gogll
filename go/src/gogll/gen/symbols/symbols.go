package symbols

import (
	"bytes"
	"fmt"
	"gogll/ast"
	"io/ioutil"
	"os"
	"path"
)

func Gen(dir string) {
	buf := new(bytes.Buffer)
	for _, sym := range ast.GetSymbols() {
		fmt.Fprintf(buf, "%s\n", sym)
	}
	if err := os.MkdirAll(dir, 0731); err != nil {
		fail(err)
	}
	if err := ioutil.WriteFile(path.Join(dir, "symbols.txt"), buf.Bytes(), 0731); err != nil {
		fail(err)
	}
}

func fail(err error) {
	fmt.Printf("Error writing symbols file: %s\n", err)
	os.Exit(1)
}
