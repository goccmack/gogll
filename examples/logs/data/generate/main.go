package main

import (
	"bytes"

	"github.com/goccmack/goutil/ioutil"
)

const data = `hello.world.com:443 10.11.12.13 - howdee3 [16/Jun/2020:09:23:14 +0200] "GET /blah/blah/blah/blah/index.html HTTP/1.1" 204 385 "https://hello.world.com/blah/blah/blah/blah/index.html" "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:77.0)"` + "\n"

func main() {
	w := new(bytes.Buffer)
	for i := 0; i < 100_000; i++ {
		w.WriteString(data)
	}
	if err := ioutil.WriteFile("../test.log", w.Bytes()); err != nil {
		panic(err)
	}
}
