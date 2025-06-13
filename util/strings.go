package util

import (
	"regexp"
	"strconv"
	"strings"
	"unicode"
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

func ChineseCount(str1 string) (count int) {
	for _, char := range str1 {
		if unicode.Is(unicode.Han, char) {
			count++
		}
	}
	return
}

//func ShowLen(str1 string) int {
//	n := 0
//	for _, char := range str1 {
//		if len(string(char)) > 2 {
//			n += 1
//		} else {
//			n++
//		}
//	}
//	return n
//}

func FirstChineseIndex(str1 string) int {
	index := -1
	for _, char := range str1 {
		index++
		if unicode.Is(unicode.Han, char) {
			return index
		}
	}
	return -1
}

func ChineseAlignNum(str1 string, num int) int {
	count := ChineseCount(str1)
	if count > num {
		return num
	}
	return num - count
}

func ChineseAlignPattern(str1 string, num int) string {
	alignNum := ChineseAlignNum(str1, num)
	return "\n%" + strconv.Itoa(alignNum) + "s"
}
