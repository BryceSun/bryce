package util

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

func Clear(s string) string {
	s = strings.TrimSpace(s)
	s = RemoveSpaceLine(s)
	return RemoveTailSpace(s)
}

func RemoveSpaceLine(s string) string {
	if SpaceLineReg.MatchString(s) {
		return SpaceLineReg.ReplaceAllString(s, "\n")
	}
	return s
}

func RemoveTailSpace(s string) string {
	if TailSpaceReg.MatchString(s) {
		return TailSpaceReg.ReplaceAllString(s, "\n")
	}
	return s
}

func Expand(s string, m map[string]string) string {
	for k, v := range m {
		s = strings.ReplaceAll(s, "${"+k+"}", v)
	}
	return s
}
