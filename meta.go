package main

import (
	"example.com/bryce/dissolve"
	"example.com/bryce/quiz"
)

type TextBlock struct {
	tittle    string      `quiz:"hide |下面进入${.}部分"`
	code      string      `quiz:"hide |代码:"`
	attention []string    `quiz:"check|${tittle}第${i}个口诀:"`
	statement string      `quiz:"show |内容"`
	questions []*Question `quiz:"show |第${i}个口诀:"`
	prev      *TextBlock
	subBlocks []*TextBlock
}

type Question struct {
	topic       string
	answer      string `quiz:"check  |${topic}的答案是"`
	explanation string `quiz:"show"`
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

func (t *TextBlock) Tittle() string {
	return t.tittle
}

func (t *TextBlock) SetTittle(s string) {
	t.tittle = s
}

func (t *TextBlock) Content() string {
	return t.statement
}

func (t *TextBlock) SetContent(s string) {
	t.statement = s
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
