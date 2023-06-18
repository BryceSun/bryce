package main

import "flag"

func main() {
	flag.Func("scan", "used to scan file", scanPlus)
	flag.Func("test", "used to test", test)
	flag.Parse()
}
