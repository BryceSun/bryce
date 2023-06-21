package main

import (
	"example.com/bryce/quiz"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
)

const (
	Ext             = ".md"
	CodeDelim       = "```"
	TagDelim        = "--"
	BoldDelim       = "**"
	KeyDelim        = "=="
	HeaderPrefixExp = ` *[-#*]+ (\*\*)?`
)

var HeaderPrefixReg = regexp.MustCompile(HeaderPrefixExp)

// 根据文件路径名扫描文档
func scanPlus(s string) error {
	_, err := scanDocument(s)
	if err != nil {
		return err
	}
	return nil
}

// 基于树形分析将词条和内容分拣后再以自定义处理器分别处理两者
func scanDocument(s string) (*TextBlock, error) {
	if Ext != path.Ext(s) {
		return nil, fmt.Errorf("ext of the file must be md")
	}
	content, err := os.ReadFile(s)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	_, file := filepath.Split(s)
	file = file[:strings.Index(file, ".")]
	document := newDocument(file)
	//先格式化，去除冗余字符
	err = quiz.Parse(string(content), document)
	return document, nil
}

// 设置自定义处理函数，和获取缩进表达式的函数
func newDocument(fileName string) *TextBlock {
	document := new(TextBlock)
	document.name = fileName
	document.handleFunc = func(tree quiz.Tree, indent string) error {
		element := tree.(*TextBlock)
		parseName(element, indent)
		parseStatement(element, indent)
		element.setHasQuiz()
		return nil
	}
	document.indentRegFunc = func() *regexp.Regexp { return HeaderPrefixReg }
	return document
}

func parseStatement(tb *TextBlock, indent string) {
	if strings.Contains(indent, tb.statement) {
		tb.statement = ""
	}
	statement := strings.TrimSpace(tb.statement)
	if strings.HasPrefix(tb.statement, CodeDelim) {
		tb.code = strings.Trim(statement, CodeDelim)
	}
}

func parseName(tb *TextBlock, indent string) {
	log.Println("解析前词条：", tb.name, tb.tag)
	log.Println("解析前缩进：", indent)
	name := strings.TrimPrefix(tb.name, indent)
	if strings.Contains(indent, BoldDelim) {
		name = strings.TrimSuffix(name, BoldDelim)
		tb.name = name
	} else if strings.Contains(name, KeyDelim) {
		name = strings.Trim(name, KeyDelim)
		prev := tb.prev
		prev.attention = append(prev.attention, name)
		prev.setHasQuiz()
		prev.removeTree(tb)
		return
	}
	tb.name = name
	if strings.Contains(name, TagDelim) {
		entry := strings.Split(name, TagDelim)
		tb.name, tb.tag = entry[0], entry[1]
	}
	log.Println("解析后词条：", tb.name, tb.tag)
}
