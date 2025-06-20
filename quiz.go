package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"strconv"
	"strings"
	"time"

	"example.com/bryce/quiz"
	"example.com/bryce/util"
)

const (
	Prefix  = "                                                     "
	Welcome = Prefix + "欢迎使用子匀问答!\n" +
		Prefix + "统计本次得分请输入S\n" +
		Prefix + "退出当前程序请输入Q\n" +
		Prefix + "跳过当前问答请输入K\n" +
		Prefix + "跳至下一题组请输入KK"
	GoodBye = Prefix + "感谢你的使用！\n" +
		Prefix + "期待下次更好的你！\n" +
		Prefix + "祝你生活愉快！"
)

var (
	Encouragement = [6]string{
		"再想想,你可以的!",
		"我不信你只有这种程度,再来!",
		"加油,坚持就是胜利!!!",
		"再做不出来就别吃饭了!",
		"这都错了啊,不是很简单的吗?",
		"答成这样可是要挨打的啊!",
	}
	Praise = [6]string{
		"不愧是你,这都做出来了!",
		"哈哈哈,不用扣鸡腿了!",
		"哟,不错嘛,又拿了一血!",
		"666!在下佩服!",
		"我不信你还能过关,再来!",
		"你这正确率,放在过去那就是状元!",
	}
	Prt = util.NewPrinter("", "%s>| %s")
)

func showWith(doc *TextBlock) {
	log.SetOutput(io.Discard)
	engine := quiz.NewTextEngine(doc)
	engine.RegisterGuardFilter(showWelcome)
	engine.RegisterGuardFilter(showWrongEntrise)
	engine.RegisterEntryFilter(printEmptyLine)
	//engine.RegisterEntryFilter(showSpendedTime)
	engine.RegisterEntryFilter(setWrongEntrise)
	engine.RegisterCoreFilter(setPath)
	engine.RegisterOrder("K", skip)
	engine.RegisterOrder("KK", skip2)
	engine.RegisterOrder("KS", searchAndSkip)
	engine.RegisterOrder("Q", skipToHead)
	engine.RegisterOrder(quiz.EncourageFunKey, encourage)
	engine.RegisterOrder(quiz.PraiseFunKey, praise)
	engine.RegisterOrder(quiz.TittleFunKey, printTittle)
	engine.RegisterOrder(quiz.StateFunKey, printStatement)
	err := engine.Start()
	if err != nil {
		log.Println(err)
	}
}
func showWelcome(e *quiz.TextEngine) error {
	fmt.Println(Welcome)
	fmt.Println()
	if err := e.Start(); err != nil {
		log.Println(err)
		return err
	}
	fmt.Println()
	fmt.Println(GoodBye)
	return nil
}

func skip(e *quiz.TextEngine) (err error) {
	n := 1
	if len(e.Input()) > 2 {
		n, err = strconv.Atoi(e.Input()[2])
		if err != nil {
			return nil
		}
	}
	e.SetSkipN(n)
	entry := e.CurrentEntry()
	Prt.Printf("%s%s", entry.Tittle, entry.Content)
	util.Scanln()
	return
}

func searchAndSkip(e *quiz.TextEngine) (err error) {
	if len(e.Input()) < 3 {
		return nil
	}
	s := e.Input()[2]
	s = strings.TrimSpace(s)
	if s == "" {
		return
	}
	e.LocalToSpecificSection(s)
	return nil
}

func skip2(e *quiz.TextEngine) error {
	if !e.LocateToNextSection() {
		Prt.Println("this action cannot be used at this point")
	}
	return nil
}

func skipToHead(e *quiz.TextEngine) error {
	e.LocateTo(e.HeadText)
	fmt.Println()
	return nil
}

func printEmptyLine(e *quiz.TextEngine) error {
	if err := e.CheckEntry(); err != nil {
		return err
	}
	fmt.Println()
	return nil
}
func showSpendedTime(e *quiz.TextEngine) error {
	start := time.Now()
	timeKey := "startTime"
	t := e.UserCache[timeKey]
	if t != nil {
		start = t.(time.Time)
	}
	if err := e.CheckEntry(); err != nil {
		return err
	}
	switch {
	case e.Right:
		Prt.Printf("本次作答用时:%d秒\n", int(time.Since(start).Seconds()))
		delete(e.UserCache, timeKey)
		return nil
	case e.HasSkip():
		delete(e.UserCache, timeKey)
		return nil
	default:
		e.UserCache[timeKey] = start
		return nil
	}
}

func encourage(e *quiz.TextEngine) error {
	//Prt.Println(Encouragement[rand.Intn(6)])
	Prt.Println("回答错误！回答错误！回答错误！")
	return nil
}

func praise(e *quiz.TextEngine) error {
	//Prt.Println(Praise[rand.Intn(6)])
	Prt.Println("恭喜你答对了！")
	return nil
}

func printTittle(e *quiz.TextEngine) error {
	entry := e.CurrentEntry()
	Prt.Print(entry.Tittle)
	return nil
}

func printStatement(e *quiz.TextEngine) (err error) {
	content := e.CurrentEntry().Content
	reader := bufio.NewReader(strings.NewReader(content))
	readString := ""
	for err == nil {
		readString, err = reader.ReadString('\n')
		fmt.Print(Prefix + readString)
	}
	if e.CurrentEntry().Kind != quiz.Confirm {
		fmt.Println()
		fmt.Println()
	} else {
		fmt.Printf("\n%s", Prefix)
	}
	log.Println(e)
	return nil
}

func setPath(e *quiz.TextEngine) error {
	text := e.CurrentText.(*TextBlock)
	if text.Tittle == "" {
		return nil
	}
	pathKey := "tittles"
	var path []string
	p := e.UserCache[pathKey]
	if p != nil {
		path = p.([]string)
		for i, s := range path {
			if s == text.prev.Tittle {
				path = path[:i+1]
				break
			}
		}
	}
	path = append(path, text.Tittle)
	e.UserCache[pathKey] = path
	Prt.Prefix = Prefix + strings.Join(path, "> ")
	return nil
}

func setWrongEntrise(e *quiz.TextEngine) error {
	if err := e.CheckEntry(); err != nil {
		return err
	}
	if e.Right || e.HasSkip() {
		return nil
	}
	key := "wrongentrise"
	wrongEntrise := util.Load[map[*quiz.EntryQuiz]string](e.UserCache, key)
	if wrongEntrise == nil {
		wrongEntrise = map[*quiz.EntryQuiz]string{}
	}
	wrongEntrise[e.CurrentEntry()] = Prt.Prefix
	e.UserCache[key] = wrongEntrise
	return nil
}

func showWrongEntrise(e *quiz.TextEngine) error {
	if err := e.Start(); err != nil {
		return err
	}
	if e.HasSkip() {
		return nil
	}
	key := "wrongentrise"
	wrongEntrise := util.Load[map[*quiz.EntryQuiz]string](e.UserCache, key)
	if wrongEntrise == nil {
		return nil
	}
	entrise := util.Keys[*quiz.EntryQuiz](wrongEntrise)
	if len(entrise) == 0 {
		return nil
	}
	fmt.Println(Prefix + "	<下面进入纠错模式>")
	e.SetQuizEntrys(entrise)
	//在展示单个词条前设置打印前缀
	setPrefix := func(e *quiz.TextEngine) error {
		entry := e.CurrentEntry()
		Prt.Prefix = wrongEntrise[entry]
		return e.CheckEntry()
	}
	e.RegisterEntryFilter(setPrefix)
	e.ShowQuizEntrys()
	fmt.Println(Prefix + "	<纠错模式结束>")
	return nil
}
