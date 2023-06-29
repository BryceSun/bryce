package main

import (
	"example.com/bryce/dissolve"
	"example.com/bryce/quiz"
)

type TextBlock struct {
	Tittle    string      `json:"Tittle" quiz:"head |下面进入<${Tittle}>部分"` // head用于表示只展示标签定义的内容
	Code      string      `json:"Code" quiz:"show"`
	Attention []string    `json:"Attention" quiz:"check|<${Tittle}>第${i}个口诀:"`
	Statement string      `json:"Statement" quiz:"show |内容"`
	Questions []*Question `json:"Questions" quiz:"show |那么第${i}个问题"`
	prev      *TextBlock
	subBlocks []*TextBlock
}

type Question struct {
	Topic       string
	Answer      string `quiz:"check  |\"${Topic}\"的答案是："`
	Explanation string `quiz:"show"`
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

func (t *TextBlock) GetTittle() string {
	return t.Tittle
}

func (t *TextBlock) SetTittle(s string) {
	t.Tittle = s
}

func (t *TextBlock) Content() string {
	return t.Statement
}

func (t *TextBlock) SetContent(s string) {
	t.Statement = s
}

func (t *TextBlock) NewTree() dissolve.Tree {
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
