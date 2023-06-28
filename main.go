package main

import (
	"flag"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	flag.Func("scan", "used to scan file", scanPlus)
	flag.Func("test", "used to test", test)
	flag.Func("store", "used to store content of a file to database", store)
	flag.Func("load", "used to show and test content from database", load)
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

func store(s string) error {
	document, err := scan(s)
	if err != nil {
		return err
	}
	if document != nil {
		return storeWithDB(document)
	}
	return nil
}

func load(s string) error {
	document, err := LoadFromDB(s)
	if err != nil {
		return err
	}
	if document != nil {
		showWith(document)
	}
	return nil
}
