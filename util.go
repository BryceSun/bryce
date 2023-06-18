package main

import (
	"regexp"
	"strings"
)

const (
	SpaceLineExp = `\n\s*\n`
	TailSpaceExp = `\s+\n`
)

var (
	SpaceLineReg = regexp.MustCompile(SpaceLineExp)
	TailSpaceReg = regexp.MustCompile(TailSpaceExp)
)

func clear(s string) string {
	s = strings.TrimSpace(s)
	//todo 去除无意义的行
	s = removeSpaceLine(s)
	return removeTailSpace(s)
}

func removeSpaceLine(s string) string {
	if SpaceLineReg.MatchString(s) {
		return SpaceLineReg.ReplaceAllString(s, "\n")
	}
	return s
}

func removeTailSpace(s string) string {
	if TailSpaceReg.MatchString(s) {
		return TailSpaceReg.ReplaceAllString(s, "\n")
	}
	return s
}
