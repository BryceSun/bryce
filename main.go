package main

import (
	"example.com/bryce/quiz"
	"flag"
)

func main() {
	flag.Func("scan", "used to scan file", scanPlus)
	flag.Func("test", "used to test", test)
	flag.Parse()
}

func test(s string) error {
	document, err := scanDocument(s)
	if err != nil {
		return err
	}
	if document != nil {
		engine := quiz.NewTextEngine(document)
		engine.Start()
	}
	return nil
}
