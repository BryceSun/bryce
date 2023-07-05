package main

import (
	"example.com/bryce/dissolve"
	"example.com/bryce/quiz"
	"strings"
)

type TextBlock struct {
	Tittle    string `json:"Tittle" quiz:"head confirm |下面进入<${Tittle}>部分"` // head用于表示只展示标签定义的内容
	Attention string `json:"Attention" quiz:"confirm"`
	Mnemonic  string `json:"Mnemonic" quiz:"check|请输入口诀:"`
	Statement string `json:"Statement" quiz:"show"`
	Question  string `json:"Question"`
	Answer    string `json:"Answer" quiz:"check |\"${Question}\"的答案是:"`
	prev      *TextBlock
	subBlocks []*TextBlock
}

func (t *TextBlock) Prev() quiz.QText {
	return t.prev
}

func (t *TextBlock) Subs() []quiz.QText {
	var qts []quiz.QText
	for _, sub := range t.subBlocks {
		qts = append(qts, sub)
	}
	return qts
}

func (t *TextBlock) SetTittle(s string) {
	t.Tittle = s
}

func (t *TextBlock) SetContent(s string) {
	if strings.TrimSpace(s) == "" {
		return
	}
	tree := t.subTree()
	tree.Statement = strings.TrimSuffix(s, "\n\r")
}

func (t *TextBlock) NewTree() dissolve.Tree {
	return t.subTree()
}

func (t *TextBlock) subTree() *TextBlock {
	n := new(TextBlock)
	n.prev = t
	t.subBlocks = append(t.subBlocks, n)
	return n
}

func (t *TextBlock) removeTree(s *TextBlock) {
	i := 0
	for _, sub := range t.subBlocks {
		if sub != s {
			t.subBlocks[i] = sub
			i++
		}
	}
	t.subBlocks = t.subBlocks[:i]
}
