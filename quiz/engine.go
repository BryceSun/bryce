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
	hIndex      int    //过滤器索引
	locIndex    int    //当前节点在父节点中的索引
	entryIndex  int    //最近一次的词条在词条集合中的索引
	input       string //用户最近一次的输入
	headText    QText  //顶端节点
	locText     QText  //定位节点
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
		tc.ScanAndTest()
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
	tc.setQuizEntrys()
	tc.showQuizEntrys()
	for _, text := range tc.currentText.Subs() {
		if tc.locText != nil && tc.locText != text {
			continue
		}
		tc.locText = nil
		tc.currentText = text
		tc.ScanAndTest()
	}
}

// 设置题集
func (tc *TextEngine) setQuizEntrys() {
	tc.quizEntrys = parseQText(tc.currentText)
}

// LocateTo 重定位到某节点
func (tc *TextEngine) LocateTo(text QText) {
	//增加偏移量
	tc.locText = text
}

// LocateToNextText  重定位到下一小节
func (tc *TextEngine) LocateToNextText() {
	if tc.locText == nil {
		tc.locText = tc.currentText
	}
	tc.setLocIndex()
	if tc.locToRighText() {
		return
	}
	if tc.locToUpperText() && tc.locToRighText() {
		return
	}
	tc.locText = nil
}

// LocateToNextSection 重定位到下一大节
func (tc *TextEngine) LocateToNextSection() {
	if tc.locText == nil {
		tc.locText = tc.currentText
	}
	tc.setLocIndex()
	if tc.locToUpperText() && tc.locToRighText() {
		return
	}
	tc.locText = nil
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
	for i, entry := range tc.quizEntrys {
		tc.right = false
		tc.entryIndex = i
		if !entry.isTest {
			tc.excFuncOrPrintln(tittleFunKey, entry.tittle)
			if strings.TrimSpace(entry.content) != "" {
				tc.excFuncOrPrintln(stateFunKey, entry.content)
			}
			continue
		}
		for !tc.right {
			tc.hIndex = 0
			tc.CheckEntry()
			if tc.locText != nil {
				return
			}
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
	keyFunc = tc.getHandler(tc.input)
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

func (tc *TextEngine) RegisterOrder(k string, h Handler) {
	tc.userOrder[k] = h
}

func (tc *TextEngine) registerOrder(k string, h Handler) {
	tc.defaultOrder[k] = h
}

func (tc *TextEngine) RegisterFilter(h Handler) {
	tc.rHandlers = append(tc.rHandlers, h)
}

func (tc *TextEngine) SetToText(text QText) {
	tc.currentText = text
	tc.setLocIndex()
}

func (tc *TextEngine) locToUpperText() bool {
	text := tc.locText
	if reflect.ValueOf(text.Prev()).IsNil() {
		return false
	}
	tc.locText = text.Prev()
	tc.setLocIndex()
	return true
}

func (tc *TextEngine) locToRighText() bool {
	text := tc.locText
	if reflect.ValueOf(text.Prev()).IsNil() {
		return false
	}
	subTexts := text.Prev().Subs()
	if tc.locIndex == len(subTexts)-1 {
		return false
	}
	tc.locIndex++
	tc.locText = subTexts[tc.locIndex]
	return true
}

// 重定位或向父节点移动时需要调用此方法
func (tc *TextEngine) setLocIndex() {
	if tc.locText == nil {
		tc.locText = tc.currentText
	}
	text := tc.locText
	if reflect.ValueOf(text.Prev()).IsNil() {
		tc.locIndex = 0
		return
	}
	for i, subText := range text.Prev().Subs() {
		if subText == text {
			tc.locIndex = i
			return
		}
	}
}
