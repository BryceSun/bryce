package main

import (
	"flag"
)

func main() {
	flag.Func("scan", "used to scan file", scanPlus)
	flag.Func("test", "used to test", test)
	flag.Parse()
}

func test(s string) error {
	document, err := scan(s)
	if err != nil {
		return err
	}
	if document != nil {
		showWith(document)
	}
	return nil
}
