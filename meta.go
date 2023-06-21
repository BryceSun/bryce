package main

import (
	"example.com/bryce/quiz"
	"regexp"
)

type TextBlock struct {
	name          string   `quiz:"${name}的${tail}是:,head"`
	code          string   `quiz:"代码,false,hide"`
	tag           string   `quiz:"答案,true,show"`
	key           string   `quiz:"重点,false,show"`
	attention     []string `quiz:"第${i}个口诀,true,show"`
	statement     string   `quiz:"内容,false,show"`
	hasQuiz       bool
	prev          *TextBlock
	subBlocks     []*TextBlock
	indentRegFunc func() *regexp.Regexp
	handleFunc    func(QuizTree, string) error
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

func (t *TextBlock) IndentReg() *regexp.Regexp {
	return t.indentRegFunc()
}

func (t *TextBlock) Tittle() string {
	return t.name
}

func (t *TextBlock) SetTittle(s string) {
	t.name = s
}

func (t *TextBlock) Content() string {
	return t.statement
}

func (t *TextBlock) SetContent(s string) {
	t.statement = s
}

func (t *TextBlock) NewTree() QuizTree {
	n := new(TextBlock)
	t.appendSub(n)
	return n
}

func (t *TextBlock) removeTree(tree QuizTree) {
	i := 0
	element := tree.(*TextBlock)
	for _, sub := range t.subBlocks {
		if sub != element {
			t.subBlocks[i] = sub
			i++
		}
	}
	t.subBlocks = t.subBlocks[:i]
}

func (t *TextBlock) handle() func(QuizTree, string) error {
	return t.handleFunc
}

func (t *TextBlock) appendSub(sub *TextBlock) {
	sub.prev = t
	sub.handleFunc = t.handleFunc
	sub.indentRegFunc = t.indentRegFunc
	t.subBlocks = append(t.subBlocks, sub)
}

func (t *TextBlock) setHasQuiz() {
	if len(t.tag) > 0 || len(t.attention) > 0 {
		t.hasQuiz = true
	}
}
