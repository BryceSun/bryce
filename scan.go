package main

import (
	"example.com/bryce/dissolve"
	"log"
	"regexp"
	"strings"
)

const (
	Ext               = ".md"
	CodeDelim         = "```"
	TagDelim          = "--"
	BoldDelim         = "**"
	KeyDelim          = "=="
	HeaderPrefixExp   = `#+ `
	ListBoldPrefixExp = ` *- +\*\*`
	ListPrefixExp     = ` *- +`
)

var HeaderPrefixReg = regexp.MustCompile(HeaderPrefixExp)
var BoldPrefixReg = regexp.MustCompile(ListBoldPrefixExp)
var ListPrefixReg = regexp.MustCompile(ListPrefixExp)

// 根据文件路径名扫描文档
// 基于树形分析将词条和内容分拣后再以自定义处理器分别处理两者
func scanPlus(s string) error {
	if !strings.HasSuffix(s, Ext) {
		log.Panic("必须是md文件")
	}
	t := new(TextBlock)
	dissolve.IndentRegs = []*regexp.Regexp{HeaderPrefixReg, BoldPrefixReg, ListPrefixReg}
	err := dissolve.ParseFile(s, t)
	if err != nil {
		return err
	}
	fixTextBlock(t)
	return nil
}

func scan(s string) (*TextBlock, error) {
	if !strings.HasSuffix(s, Ext) {
		log.Panic("必须是md文件")
	}
	t := new(TextBlock)
	dissolve.IndentRegs = []*regexp.Regexp{HeaderPrefixReg, BoldPrefixReg, ListPrefixReg}
	err := dissolve.ParseFile(s, t)
	if err != nil {
		return nil, err
	}
	handleTextBlock(t)
	return t, err
}

func handleTextBlock(block *TextBlock) {
	var blocks []*TextBlock
	blocks = append(blocks, block.subBlocks...)
	for _, b := range blocks {
		handleTextBlock(b)
	}
	fixTextBlock(block)
}

func fixTextBlock(block *TextBlock) {
	tittle, content := block.Tittle, block.Statement
	if content == "-" {
		content = ""
	}
	switch {
	case strings.HasSuffix(tittle, BoldDelim):
		block.Tittle = strings.TrimSuffix(tittle, BoldDelim)
		block.Statement = content
	case strings.HasPrefix(tittle, KeyDelim):
		tittle = strings.Trim(tittle, KeyDelim)
		prev := block.prev
		prev.Attention = append(prev.Attention, tittle)
		prev.removeTree(block)
	case strings.Contains(tittle, TagDelim):
		entry := strings.Split(tittle, TagDelim)
		q := &Question{
			Topic:       strings.TrimSpace(entry[0]),
			Answer:      strings.TrimSpace(entry[1]),
			Explanation: strings.TrimSpace(content),
		}
		prev := block.prev
		prev.Questions = append(prev.Questions, q)
		prev.removeTree(block)
	}
	if strings.HasPrefix(content, CodeDelim) {
		block.Code = strings.Trim(content, CodeDelim)
		block.Statement = ""
	}
}
