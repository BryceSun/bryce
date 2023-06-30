package util

import (
	"bufio"
	"log"
	"os"
	"strings"
)

func Scanln() string {
	in, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil {
		log.Println(err)
	}
	return strings.TrimSpace(in)
}
