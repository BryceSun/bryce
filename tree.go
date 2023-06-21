package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
)

var logger = log.Default()

type Tree interface {
	indentReg() *regexp.Regexp
	tittle() string
	SetTittle(string)
	content() string
	setContent(string)
	newTree() Tree
	removeTree(Tree)
	handle() func(Tree, string) error
}

type textCube struct {
	indent string
	tittle string
	block  string
}

func parse(text string, t Tree) error {
	logFile, err := os.Create("./tree_log.txt")
	if err == nil {
		logger.SetOutput(logFile)
	}
	log.Println(err)
	defer func() {
		if logFile != nil {
			err := logFile.Close()
			if err != nil {
				panic(err)
			}
		}
	}()
	text = clear(text)
	return handle(text, t)
}

func handle(text string, t Tree) error {
	reg := t.indentReg()
	subTexts := splitText(text, reg)
	for _, subText := range subTexts {
		subt := t.newTree()
		err := fillTree(subt, subText.block, subText.tittle, subText.indent)
		if err != nil {
			return err
		}
	}
	return nil
}

func fillTree(t Tree, text string, tittle string, indent string) (err error) {
	defer func() {
		// 返回前调用用户自定义函数处理数据
		err = t.handle()(t, indent)
	}()
	t.SetTittle(tittle)
	logger.Println("收录词条：", t.tittle())
	//没有下一行则返回
	lineEndIndex := strings.IndexAny(text, "\n\r")
	if lineEndIndex == -1 {
		return
	}
	text = text[lineEndIndex+1:]
	//判断是否有子词条
	if !t.indentReg().MatchString(text) {
		t.setContent(text)
		return
	}
	//获取子词条的索引
	sonIndentIndex := t.indentReg().FindStringIndex(text)
	//取索引前的数据内容
	t.setContent(text[:sonIndentIndex[0]])
	text = strings.TrimLeft(text[sonIndentIndex[0]:], "\n\r")
	return handle(text, t)
}

func splitText(text string, indentReg *regexp.Regexp) (subs []*textCube) {
	logger.Println("初始正则表达式：", indentReg)
	if len(text) == 0 {
		return
	}
	if !indentReg.MatchString(text) {
		logger.Panicln("文本不能被初始缩进正则表达式匹配")
		return
	}
	firstindent := indentReg.FindString(text)
	patternName := "indent"
	indentReg = NewIndentReg(firstindent, patternName)
	if !indentReg.MatchString(text) {
		logger.Panicln("文本不能被缩进正则表达式匹配, 表达式：", indentReg)
		return
	}
	indexes := indentReg.FindAllStringSubmatchIndex(text, -1)
	indentIndex := indentReg.SubexpIndex(patternName) * 2

	l := len(indexes)
	subs = make([]*textCube, len(indexes))
	start, end := indexes[0][indentIndex], -1
	for i := 0; i < l; i++ {
		//获取词条最大范围的文本
		j := i + 1
		subText := ""
		if j == l {
			subText = text[start:]
		} else {
			end = indexes[j][indentIndex]
			subText = text[start:end]
		}
		t := new(textCube)
		subs[i] = t
		t.block = strings.Trim(subText, "\n\r")
		start = end
		//获取词条
		start2 := indexes[i][0]
		end2 := indexes[i][1]
		t.tittle = strings.Trim(text[start2:end2], "\n\r")
		logger.Printf("分割到第%d个词条%s", i, subs[i])
		//获取词条缩进
		start2 = indexes[i][indentIndex]
		end2 = indexes[i][indentIndex+1]
		t.indent = strings.Trim(text[start2:end2], "\n\r")
	}
	return
}

func NewIndentReg(indent, patternName string) *regexp.Regexp {
	pattern := fmt.Sprintf(`\n?^?(?P<%s>%s).*`, patternName, regexp.QuoteMeta(indent))
	indentReg, err := regexp.Compile(pattern)
	if err != nil {
		logger.Panicln(err)
	}
	return indentReg
}
