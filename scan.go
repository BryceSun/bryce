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
	NumberPrefixExp   = `[0-9]+\. `
)

var HeaderPrefixReg = regexp.MustCompile(HeaderPrefixExp)
var BoldPrefixReg = regexp.MustCompile(ListBoldPrefixExp)
var ListPrefixReg = regexp.MustCompile(ListPrefixExp)
var NumberPrefixReg = regexp.MustCompile(NumberPrefixExp)

// 根据文件路径名扫描文档
// 基于树形分析将词条和内容分拣后再以自定义处理器分别处理两者
func scanPlus(s string) error {
	_, err := scan(s)
	if err != nil {
		return err
	}
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
	for _, b := range block.subBlocks {
		handleTextBlock(b)
	}
	handleTittle(block)
	handleStatement(block)
}

func handleTittle(block *TextBlock) {
	tittle := block.Tittle
	if strings.TrimSpace(tittle) == "" {
		return
	}
	block.Tittle = ""
	switch {
	case strings.HasSuffix(tittle, BoldDelim):
		block.Tittle = strings.TrimSuffix(tittle, BoldDelim)
		return
	case strings.HasPrefix(tittle, KeyDelim):
		tittle = strings.Trim(tittle, KeyDelim)
		block.Mnemonic = tittle
		return

	case strings.Contains(tittle, TagDelim):
		entry := strings.Split(tittle, TagDelim)
		block.Question = strings.TrimSpace(entry[0])
		block.Answer = strings.TrimSpace(entry[1])
		return
	}
	block.Tittle = tittle
}

func handleStatement(block *TextBlock) {
	content := block.Statement
	if strings.TrimSpace(content) == "" {
		return
	}
	block.Statement = ""
	if strings.HasPrefix(content, CodeDelim) {
		block.Statement = strings.TrimSpace(content)
		block.Statement = strings.Trim(block.Statement, CodeDelim)
		block.Statement = strings.TrimSpace(block.Statement)
		return
	}
	if NumberPrefixReg.MatchString(content) {
		block.Attention = content
		return
	}
	if content != "-" {
		block.Statement = content
	}
}
