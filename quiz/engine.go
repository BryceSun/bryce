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
	TittleFunKey    = "printTittle"
	StateFunKey     = "printState"
	PraiseFunKey    = "printPraise"
	EncourageFunKey = "PrintEncourage"
	GoodByeFunKey   = "PrintGoodbye"
	AtferSetEntry   = "atferSetEntry"
)

type Handler func(ctx *TextEngine) error

type QText interface {
	Prev() QText
	Subs() []QText
}

type TextEngine struct {
	Right       bool   //用户输入是否正确
	hIndex      int    //运行前后的过滤器索引
	eIndex      int    //测试题(词条)过滤器索引
	locIndex    int    //目标节点在父节点中的索引,定位时会设置这个字段
	entryIndex  int    //最近一次的词条在词条集合中的索引
	offset      int    //被忽略输出的词条的次数,即下一个要打印的词条距当下词条的偏移量
	input       string //用户最近一次的输入
	HeadText    QText  //顶端节点
	locText     QText  //定位节点
	CurrentText QText
	quizEntrys  []*EntryQuiz //词条集合
	rHandlers   []Handler    //作用于运行前后的过滤器集合
	eHandlers   []Handler    //作用于测试题前后的过滤器集合
	//once         sync.Once
	defaultOrder map[string]Handler //保存系统默认插件
	userOrder    map[string]Handler //保存用户自定义插件
	UserCache    map[string]any     //保存用户自定义数据的内存空间
}

func (tc *TextEngine) Start() {
	for tc.hIndex < len(tc.rHandlers) {
		if tc.hIndex == -1 {
			return
		}
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
		CurrentText:  qt,
		HeadText:     qt,
		rHandlers:    []Handler{},
		eHandlers:    []Handler{},
		defaultOrder: map[string]Handler{},
		userOrder:    map[string]Handler{},
		UserCache:    map[string]any{},
	}
	return engine
}

func (tc *TextEngine) ScanAndTest() {
	tc.setQuizEntrys()
	f := tc.getHandler(AtferSetEntry)
	if f != nil {
		_ = f(tc)
	}
	tc.showQuizEntrys()
	for _, text := range tc.CurrentText.Subs() {
		if tc.locText != nil && tc.locText != text {
			continue
		}
		tc.locText = nil
		tc.CurrentText = text
		tc.ScanAndTest()
	}
}

// 设置题集
func (tc *TextEngine) setQuizEntrys() {
	tc.quizEntrys = parseQText(tc.CurrentText)
}

// LocateTo 重定位到某节点
func (tc *TextEngine) LocateTo(text QText) {
	//增加偏移量
	tc.locText = text
}

// LocateToNextText  重定位到下一小节
func (tc *TextEngine) LocateToNextText() bool {
	if tc.locText == nil {
		tc.locText = tc.CurrentText
	}
	tc.setLocIndex()
	if tc.locToRighText() {
		return true
	}
	if tc.locToUpperText() && tc.locToRighText() {
		return true
	}
	tc.locText = nil
	return false
}

// LocateToNextSection 重定位到下一大节
func (tc *TextEngine) LocateToNextSection() bool {
	if tc.locText == nil {
		tc.locText = tc.CurrentText
	}
	tc.setLocIndex()
	if tc.locToUpperText() && tc.locToRighText() {
		return true
	}
	tc.locText = nil
	return false
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
	tc.initEntryRange()
	for tc.hasNextEntry() {
		if tc.HasSkip() {
			tc.skip()
			continue
		}
		entry := tc.nextEntry()
		tc.Right = false
		if !entry.IsTest {
			tc.excFuncOrPrintln(TittleFunKey, entry.Tittle)
			if strings.TrimSpace(entry.Content) != "" {
				tc.excFuncOrPrintln(StateFunKey, entry.Content)
			}
			continue
		}
		for !tc.Right {
			tc.eIndex = 0
			tc.CheckEntry()
			if tc.locText != nil {
				return
			}
			if tc.HasSkip() {
				tc.skip()
				break
			}
		}
	}
}

// CheckEntry 展示测试题目的前后调用拦截器
func (tc *TextEngine) CheckEntry() {
	for tc.eIndex < len(tc.eHandlers) {
		if tc.eIndex == -1 {
			return
		}
		tc.eIndex++
		err := tc.eHandlers[tc.eIndex-1](tc)
		if err != nil {
			panic(err)
		}
	}
	if tc.eIndex != -1 {
		tc.scanAndCheck()
		tc.eIndex = -1
	}
}

func (tc *TextEngine) scanAndCheck() {
	entry := tc.CurrentEntry()
	answer := entry.Content
	tc.excFuncOrPrintln(TittleFunKey, entry.Tittle)
	var keyFunc Handler
	if tc.getUserInput() && tc.input == answer {
		tc.Right = true
		tc.excFuncOrPrintln(PraiseFunKey, "回答正确！")
		return
	}
	keyFunc = tc.getHandler(tc.input)
	if keyFunc == nil {
		tc.excFuncOrPrintln(EncourageFunKey, "回答错误！")
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

func (tc *TextEngine) RegisterGuardFilter(h Handler) {
	tc.rHandlers = append(tc.rHandlers, h)
}

func (tc *TextEngine) RegisterEntryFilter(h Handler) {
	log.Println("append filter...")
	tc.eHandlers = append(tc.eHandlers, h)
}

func (tc *TextEngine) SetToText(text QText) {
	tc.CurrentText = text
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
		tc.locText = tc.CurrentText
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

func (tc *TextEngine) HasSkip() bool {
	return tc.offset > 0
}

func (tc *TextEngine) HasLocate() bool {
	return tc.locText != nil
}

// SkipEntryN  跳过N个词条
func (tc *TextEngine) skip() {
	tc.offset--
}

// SetSkipOnce 路过当前词条
func (tc *TextEngine) SetSkipOnce() {
	tc.offset = 1
}

// SetSkipN 跳过N个词条
func (tc *TextEngine) SetSkipN(n int) {
	tc.offset = n
}

// 初始化词条迭代状态,需要配合 hasNextEntry 方法使用
func (tc *TextEngine) initEntryRange() {
	tc.entryIndex = -1
}

// 是否有下一个词条,需要配合 initEntryRange 和 nextEntry 方法使用
func (tc *TextEngine) hasNextEntry() bool {
	b := tc.entryIndex < len(tc.quizEntrys)-1
	if !b {
		tc.entryIndex = 0
	}
	return b
}

// 取下一个词条,需要配合 hasNextEntry 方法使用
func (tc *TextEngine) nextEntry() *EntryQuiz {
	tc.entryIndex++
	return tc.CurrentEntry()
}

// CurrentEntry 取当前词条
func (tc *TextEngine) CurrentEntry() *EntryQuiz {
	return tc.quizEntrys[tc.entryIndex]
}
