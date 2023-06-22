package main

import (
	"example.com/bryce/dissolve"
	"log"
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
// 基于树形分析将词条和内容分拣后再以自定义处理器分别处理两者
func scanPlus(s string) error {
	if !strings.HasSuffix(s, Ext) {
		log.Panic("必须是md文件")
	}
	t := new(TextBlock)
	dissolve.IndentReg = HeaderPrefixReg
	dissolve.FillTree = fillTextBlock
	return dissolve.ParseFile(s, t)
}

func scan(s string) (*TextBlock, error) {
	if !strings.HasSuffix(s, Ext) {
		log.Panic("必须是md文件")
	}
	t := new(TextBlock)
	dissolve.IndentReg = HeaderPrefixReg
	dissolve.FillTree = fillTextBlock
	err := dissolve.ParseFile(s, t)
	return t, err
}

func fillTextBlock(tree dissolve.Tree, cube *dissolve.TextCube) dissolve.Tree {
	indent, tittle, content := cube.Indent, cube.Tittle, cube.Content
	textBlock := tree.(*TextBlock)
	switch {
	case strings.Contains(indent, BoldDelim):
		textBlock.tittle = strings.TrimSuffix(tittle, BoldDelim)
		textBlock.statement = content
		return textBlock

	case strings.Contains(tittle, KeyDelim):
		tittle = strings.Trim(tittle, KeyDelim)
		prev := textBlock.prev
		prev.attention = append(prev.attention, tittle)
		prev.removeTree(textBlock)
		return prev

	case strings.Contains(tittle, TagDelim):
		entry := strings.Split(tittle, TagDelim)
		q := &Question{entry[0], entry[1], content}
		prev := textBlock.prev
		prev.questions = append(prev.questions, q)
		prev.removeTree(textBlock)
		return prev
	}
	textBlock.tittle = tittle
	textBlock.statement = content
	if strings.HasPrefix(content, CodeDelim) {
		textBlock.code = strings.Trim(content, CodeDelim)
		textBlock.statement = ""
	}
	return textBlock
}
