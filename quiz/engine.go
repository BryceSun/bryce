package quiz

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"
)

// 插件函数的注册名
var (
	WelcomeFunKey   = "printWelcome"
	tittleFunKey    = "printTittle"
	stateFunKey     = "printState"
	praiseFunKey    = "printPraise"
	encourageFunKey = "PrintEncourage"
	GoodByeFunKey   = "PrintGoodbye"
)

type Handler func(ctx *TextEngine) error

type QText interface {
	Prev() QText
	Subs() []QText
}

type TextEngine struct {
	right       bool   //用户输入是否正确
	offset      int    //偏移量，用于记录迭代中时发生偏移的，当调用重定位方法，偏移量会加一，从重定位方法返回时并不会减一
	hIndex      int    //过滤器索引
	textIndex   int    //当前节点在父节点中的索引
	entryIndex  int    //最近一次的词条在词条集合中的索引
	input       string //用户最近一次的输入
	headText    QText  //顶端节点
	currentText QText
	quizEntrys  []*EntryQuiz //词条集合
	rHandlers   []Handler    //作用于运行前后的过滤器集合
	eHandlers   []Handler    //作用于测试题前后的过滤器集合
	//once         sync.Once
	defaultOrder map[string]Handler //保存系统默认插件
	userOrder    map[string]Handler //保存用户自定义插件
	userCache    map[string]any     //保存用户自定义数据的内存空间
}

func (tc *TextEngine) Start() {
	for tc.hIndex < len(tc.rHandlers) {
		tc.hIndex++
		err := tc.rHandlers[tc.hIndex-1](tc)
		if err != nil {
			panic(err)
		}
	}
	//tc.once.Do(tc.ScanAndTest)
	if tc.hIndex != -1 {
		if tc.offset > 0 {
			tc.RightUpScan()
		} else {
			tc.ScanAndTest()
		}
		tc.hIndex = -1
	}
}

func NewTextEngine(qt QText) *TextEngine {
	if !reflect.ValueOf(qt.Prev()).IsNil() {
		log.Panicln("非头部节点，不可用")
	}
	engine := &TextEngine{
		currentText:  qt,
		headText:     qt,
		rHandlers:    []Handler{},
		eHandlers:    []Handler{},
		defaultOrder: map[string]Handler{},
		userOrder:    map[string]Handler{},
		userCache:    map[string]any{},
	}
	return engine
}

func (tc *TextEngine) ScanAndTest() {
	tc.offset++
	//记录当前偏移量
	offset := tc.offset
	defer func() { tc.offset-- }()
	tc.setQuizEntrys()
	tc.showQuizEntrys()
	for _, text := range tc.currentText.Subs() {
		tc.currentText = text
		tc.ScanAndTest()
		if offset != tc.offset {
			return
		}
	}
}

// 设置题集
func (tc *TextEngine) setQuizEntrys() {
	tc.quizEntrys = parseQText(tc.currentText)
}

// LocateTo 重定位到某节点
func (tc *TextEngine) LocateTo(text QText) {
	//增加偏移量
	tc.offset++
	tc.SetToText(text)
	tc.RightUpScan()
}

func (tc *TextEngine) RightUpScan() {
	tc.ScanAndTest()
	for tc.SetToRighText() {
		tc.ScanAndTest()
	}
	if tc.SetToUpperText() && tc.SetToRighText() {
		tc.RightUpScan()
	}
}

func showWelcome(tc *TextEngine) {
	tc.excFuncOrPrintln(WelcomeFunKey, "欢迎使用子匀问答！")
}

func showGoodBye(tc *TextEngine) {
	tc.excFuncOrPrintln(GoodByeFunKey, "欢迎下次再来！")
}

func (tc *TextEngine) excFuncOrPrintln(funcKey string, s string) {
	keyFunc := tc.getHandler(funcKey)
	if keyFunc != nil {
		err := keyFunc(tc)
		if err == nil {
			return
		}
		log.Println(err)
	}
	fmt.Println(s)
}

func (tc *TextEngine) showQuizEntrys() {
	offset := tc.offset
	for i, entry := range tc.quizEntrys {
		tc.right = false
		tc.entryIndex = i
		if !entry.isTest {
			tc.excFuncOrPrintln(tittleFunKey, entry.tittle)
			tc.excFuncOrPrintln(stateFunKey, entry.content)
			continue
		}
		for !tc.right && offset == tc.offset {
			tc.hIndex = 0
			tc.CheckEntry()
		}
	}
}

// CheckEntry 展示测试题目前调用拦截器
func (tc *TextEngine) CheckEntry() {
	for tc.hIndex < len(tc.eHandlers) {
		tc.hIndex++
		err := tc.eHandlers[tc.hIndex-1](tc)
		if err != nil {
			panic(err)
		}
	}
	if tc.hIndex != -1 {
		tc.scanAndCheck()
		tc.hIndex = -1
	}
}

func (tc *TextEngine) scanAndCheck() {
	entry := tc.quizEntrys[tc.entryIndex]
	answer := entry.content
	tc.excFuncOrPrintln(tittleFunKey, entry.tittle)
	var keyFunc Handler
	if tc.getUserInput() && tc.input == answer {
		tc.right = true
		tc.excFuncOrPrintln(praiseFunKey, "回答正确！")
		return
	}
	keyFunc = tc.getHandler(answer)
	if keyFunc == nil {
		tc.excFuncOrPrintln(encourageFunKey, "回答错误！")
		return
	}
	err := keyFunc(tc)
	if err != nil {
		log.Println(err)
	}
	return
}

func (tc *TextEngine) getUserInput() bool {
	tc.input = ""
	input, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil {
		log.Println(err)
		return false
	}
	tc.input = strings.TrimSpace(input)
	return true
}

func (tc *TextEngine) getHandler(s string) Handler {
	if len(strings.TrimSpace(s)) == 0 {
		return nil
	}
	h := tc.userOrder[s]
	if h == nil {
		return tc.defaultOrder[s]
	}
	return h
}

func (tc *TextEngine) RegisterHandler(k string, h Handler) {
	tc.userOrder[k] = h
}

func (tc *TextEngine) registerHandler(k string, h Handler) {
	tc.defaultOrder[k] = h
}

func (tc *TextEngine) SetToText(text QText) {
	tc.currentText = text
	tc.setIndex()
}

func (tc *TextEngine) SetToUpperText() bool {
	text := tc.currentText
	if reflect.ValueOf(text.Prev()).IsNil() {
		return false
	}
	tc.currentText = text.Prev()
	tc.setIndex()
	return true
}

func (tc *TextEngine) SetToRighText() bool {
	text := tc.currentText
	if reflect.ValueOf(text.Prev()).IsNil() {
		return false
	}
	subTexts := text.Prev().Subs()
	if tc.textIndex == len(subTexts)-1 {
		return false
	}
	tc.textIndex++
	tc.currentText = subTexts[tc.textIndex]
	return true
}

// 重定位或向父节点移动时需要调用此方法
func (tc *TextEngine) setIndex() {
	text := tc.currentText
	if reflect.ValueOf(text.Prev()).IsNil() {
		tc.textIndex = -1
		return
	}
	for i, subText := range text.Prev().Subs() {
		if subText == text {
			tc.textIndex = i
			return
		}
	}
}
