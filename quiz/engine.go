package quiz

import (
	"example.com/bryce/util"
	"fmt"
	"log"
	"reflect"
	"strings"
)

// 插件函数的注册名
var (
	TittleFunKey    = "printTittle"
	StateFunKey     = "printState"
	PraiseFunKey    = "printPraise"
	EncourageFunKey = "PrintEncourage"
)

type Handler func(ctx *TextEngine) error

type QText interface {
	Prev() QText
	Subs() []QText
}

type TextEngine struct {
	Right       bool     //用户输入是否正确
	rIndex      int      //运行的过滤器索引
	eIndex      int      //测试题(词条)过滤器索引
	cIndex      int      //测试题设置的过滤器索引
	locIndex    int      //目标节点在父节点中的索引,定位时会设置这个字段
	entryIndex  int      //最近一次的词条在词条集合中的索引
	offset      int      //被忽略输出的词条的次数,即下一个要打印的词条距当下词条的偏移量
	input       []string //用户最近一次的输入
	HeadText    QText    //顶端节点
	locText     QText    //定位节点
	CurrentText QText
	quizEntrys  []*EntryQuiz //词条集合
	rHandlers   []Handler    //作用于运行前后的过滤器集合
	cHandlers   []Handler    //作用于词条集合被设置前后的过滤器集合
	eHandlers   []Handler    //作用于测试题前后的过滤器集合
	//once         sync.Once
	defaultOrder map[string]Handler //保存系统默认插件
	userOrder    map[string]Handler //保存用户自定义插件
	UserCache    map[string]any     //保存用户自定义数据的内存空间
}

func (tc *TextEngine) Input() []string {
	return tc.input
}

func (tc *TextEngine) Start() (err error) {
	return tc.throughFilters(parseAndTest, &tc.rHandlers, &tc.rIndex)(tc)
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

func parseAndTest(tc *TextEngine) (err error) {
	tc.cIndex = 0
	err = tc.ParseAndSetEntrys()
	if err != nil {
		return err
	}
	tc.ShowQuizEntrys()
	for _, text := range tc.CurrentText.Subs() {
		if tc.locText != nil && tc.locText != text {
			continue
		}
		tc.locText = nil
		tc.CurrentText = text
		if err := parseAndTest(tc); err != nil {
			log.Println(err)
			return err
		}
	}
	return err
}

// 将给定的f函数包装为附加了过滤链的函数
func (tc *TextEngine) throughFilters(f Handler, filters *[]Handler, i *int) (h Handler) {
	//缓存中存在执行过滤链的函数就直接返回
	if h = tc.defaultOrder[fmt.Sprint(i)]; h != nil {
		return h
	}
	//定义附加了过滤链的函数
	h = func(ctx *TextEngine) error {
		for *i < len(*filters) {
			if *i == -1 {
				return nil
			}
			*i++
			if err := (*filters)[*i-1](ctx); err != nil {
				return err
			}
		}
		if *i == -1 {
			return nil
		}
		*i = -1
		return f(ctx)
	}
	//将执行过滤链的函数缓存
	tc.defaultOrder[fmt.Sprint(i)] = h
	return h
}

// ParseAndSetEntrys 设置题集
func (tc *TextEngine) ParseAndSetEntrys() error {
	return tc.throughFilters(parseText, &tc.cHandlers, &tc.cIndex)(tc)
}

// 设置题集
func parseText(tc *TextEngine) error {
	tc.quizEntrys = parseQText(tc.CurrentText)
	return nil
}

// LocateTo 重定位到某节点
func (tc *TextEngine) LocateTo(text QText) {
	//增加偏移量
	tc.locText = text
}

// LocateToNextSection 重定位到下一文本块
func (tc *TextEngine) LocateToNextSection() bool {
	if tc.locText == nil {
		tc.locText = tc.CurrentText
	}
	tc.setLocIndex()
	if tc.locToUpperText() {
		if !tc.locToRighText() {
			return tc.LocateToNextSection()
		}
		return true
	}
	tc.locText = nil
	return false
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

// ShowQuizEntrys 展示词条或测试
func (tc *TextEngine) ShowQuizEntrys() {
	tc.initEntryRange()
	for tc.hasNextEntry() {
		entry := tc.nextEntry()
		if tc.HasSkip() {
			tc.skip()
			continue
		}
		tc.Right = false
		//非测试题仅进行内容展示
		if entry.Kind != Test {
			if strings.TrimSpace(entry.Tittle) != "" {
				tc.excFuncOrPrintln(TittleFunKey, entry.Tittle)
			}
			if strings.TrimSpace(entry.Content) != "" {
				tc.excFuncOrPrintln(StateFunKey, entry.Content)
			}
			if entry.Kind == Confirm {
				util.Scanln()
			}
			continue
		}
		//测试题如果没被答对或是跳过则反复展示
		for !tc.Right {
			tc.eIndex = 0
			//展示测试题目的前后调用拦截器
			err := tc.CheckEntry()
			if err != nil {
				return
			}
			if tc.HasSkip() {
				tc.skip()
				break
			}
		}
	}
}

func (tc *TextEngine) CheckEntry() (err error) {
	return tc.throughFilters(checkEntry, &tc.eHandlers, &tc.eIndex)(tc)
}

func checkEntry(tc *TextEngine) (err error) {
	entry := tc.CurrentEntry()
	var keyFunc Handler
	tc.input = tc.input[:0]
	for len(tc.input) == 0 {
		tc.excFuncOrPrintln(TittleFunKey, entry.Tittle)
		in := util.Scanln()
		if in != "" {
			tc.input = append(tc.input, in)
			tc.input = append(tc.input, strings.Split(in, " ")...)
		}
	}
	if tc.input[0] == entry.Content {
		tc.Right = true
		tc.excFuncOrPrintln(PraiseFunKey, "回答正确！")
		return
	}
	keyFunc = tc.getHandler(tc.input[1])
	if keyFunc == nil {
		tc.excFuncOrPrintln(EncourageFunKey, "回答错误！")
		return
	}
	if err = keyFunc(tc); err != nil {
		log.Println(err)
	}
	return
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

func (tc *TextEngine) RegisterCoreFilter(h Handler) {
	log.Println("append core filter...")
	tc.cHandlers = append(tc.cHandlers, h)
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

// GetIndex 重定位或向父节点移动时需要调用此方法
func (tc *TextEngine) GetIndex() int {
	if tc.locText == nil {
		tc.locText = tc.CurrentText
	}
	text := tc.locText
	if reflect.ValueOf(text.Prev()).IsNil() {
		return 0
	}
	for i, subText := range text.Prev().Subs() {
		if subText == text {
			return i
		}
	}
	return 0
}

func (tc *TextEngine) HasSkip() bool {
	return tc.offset > 0 || tc.locText != nil
}

// SkipEntryN  跳过N个词条
func (tc *TextEngine) skip() {
	tc.offset--
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

func (tc *TextEngine) SetQuizEntrys(quizEntrys []*EntryQuiz) {
	tc.quizEntrys = quizEntrys
}

// CurrentEntry 取当前词条
func (tc *TextEngine) CurrentEntry() *EntryQuiz {
	return tc.quizEntrys[tc.entryIndex]
}

func (tc *TextEngine) Save(key string, a any) {
	tc.UserCache[key] = a
}
