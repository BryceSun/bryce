package main

import "regexp"

type TextBlock struct {
	name          string
	code          string
	tag           string
	key           string
	hasQuiz       bool
	attention     []string
	statement     string
	prev          *TextBlock
	subBlocks     []*TextBlock
	indentRegFunc func() *regexp.Regexp
	handleFunc    func(Tree, string) error
}

func (t *TextBlock) indentReg() *regexp.Regexp {
	return t.indentRegFunc()
}

func (t *TextBlock) tittle() string {
	return t.name
}

func (t *TextBlock) SetTittle(s string) {
	t.name = s
}

func (t *TextBlock) content() string {
	return t.statement
}

func (t *TextBlock) setContent(s string) {
	t.statement = s
}

func (t *TextBlock) newTree() Tree {
	n := new(TextBlock)
	t.appendSub(n)
	return n
}

func (t *TextBlock) removeTree(tree Tree) {
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

func (t *TextBlock) handle() func(Tree, string) error {
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
