package main

import (
	"example.com/bryce/util"
	"flag"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

func main() {
	flag.Func("scan", "used to scan file", scanPlus)
	flag.Func("test", "used to test", test)
	flag.Func("store", "used to store content of a file to database", store)
	flag.Func("load", "used to show and test content from database", load)
	flag.Func("list", "used to list text stored with database", list)
	flag.Func("key", "used to input keymap", getKeyMap)
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

func list(string) error {
	var input string
	for {
		input = ""
		for input == "" {
			names := ListTextNames()
			for _, name := range names {
				fmt.Println(name)
			}
			fmt.Print("请选择：")
			input = util.Scanln()
		}
		if input == "Q" {
			return nil
		}
		err := load(input)
		if err != nil {
			log.Print(err)
		}
	}
}
