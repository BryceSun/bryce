package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"example.com/bryce/util"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	flag.Func("scan", "used to scan file", scanPlus)
	flag.Func("test", "used to test", test)
	flag.Func("store", "used to store content of a file to database", store)
	flag.Func("load", "used to show and test content from database", load)
	flag.Func("list", "used to list text stored with database", list)
	flag.Func("show", "used to show text line by line", showTypeEnglish)
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

func showTypeEnglish(p string) error {
	if strings.TrimSpace(p) == "" {
		p = "./"
	}
	files, err := os.ReadDir(p)
	if err != nil {
		return err
	}
	err = LoadHistory(p)
	if err != nil {
		return err
	}
	go LogHistory()
	var input string
	for {
		input = ""
		for input == "" {
			fn := ""
			for _, file := range files {
				fn = file.Name()
				util.Rprintln(fn)
			}
			util.Lprint("请选择：" + fn)
			fmt.Print("请选择：")
			input = util.Scanln()
		}
		if input == "Q" {
			return SaveHistory()
		}
		input = p + "\\" + input
		err = ShowType(input)
		if err != nil {
			log.Print(err)
		}
	}
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
