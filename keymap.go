package main

import (
	"example.com/bryce/util"
	"fmt"
	"log"
	"os"
	"strings"
)

func getKeyMap(name string) error {
	f, err := os.OpenFile(name,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Panic(err)
	}
	defer f.Close()
	desc := ""
	for {
		fmt.Print("please input desc：")
		desc = util.Scanln()
		if desc == "Q" {
			return nil
		}
		fmt.Print("Please input the keymap：")
		k := util.Scanln()
		k = strings.ReplaceAll(k, "c ", "ctrl+")
		k = strings.ReplaceAll(k, "s ", "shift+")
		k = strings.ReplaceAll(k, "a ", "alt+")
		_, err := fmt.Fprintf(f, "- %s -- %s\n", desc, k)
		if err != nil {
			return err
		}
		fmt.Printf("- %s -- %s\n", desc, k)
	}
}
