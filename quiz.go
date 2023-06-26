package main

import (
	"example.com/bryce/quiz"
	"example.com/bryce/util"
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

const (
	Prefix  = "                 "
	Welcome = Prefix + "欢迎使用子匀问答!\n" +
		Prefix + "统计本次得分请输入S\n" +
		Prefix + "退出当前程序请输入Q\n" +
		Prefix + "跳过当前问答请输入K\n" +
		Prefix + "跳至下一题组请输入KK"
	GoodBye = Prefix + "就这样算了???不要怂啊!\n" +
		Prefix + "期待下次更好的你!"
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
	Prt = util.NewPrinter("", "%s  %s")
)

func showWith(doc *TextBlock) {
	engine := quiz.NewTextEngine(doc)
	engine.RegisterGuardFilter(showWelcome)
	engine.RegisterGuardFilter(showWrongEntrise)
	engine.RegisterEntryFilter(showSpendedTime)
	engine.RegisterEntryFilter(setWrongEntrise)
	engine.RegisterCoreFilter(setPath)
	engine.RegisterOrder("K", skip)
	engine.RegisterOrder("KK", skip1)
	engine.RegisterOrder("KKK", skip2)
	engine.RegisterOrder("Q", skipToHead)
	engine.RegisterOrder(quiz.EncourageFunKey, encourage)
	engine.RegisterOrder(quiz.PraiseFunKey, praise)
	engine.RegisterOrder(quiz.TittleFunKey, printTittle)
	err := engine.Start()
	if err != nil {
		log.Println(err)
	}
}
func showWelcome(e *quiz.TextEngine) error {
	fmt.Println(Welcome)
	if err := e.Start(); err != nil {
		log.Println(err)
		return err
	}
	fmt.Println(GoodBye)
	return nil
}

func skip(e *quiz.TextEngine) (err error) {
	n := 1
	if len(e.Input()) > 2 {
		n, err = strconv.Atoi(e.Input()[2])
		if err != nil {
			return
		}
	}
	e.SetSkipN(n)
	entry := e.CurrentEntry()
	Prt.Printf("%s %s\n", entry.Tittle, entry.Content)
	return
}

func skip1(e *quiz.TextEngine) error {
	e.LocateToNextText()
	return nil
}

func skip2(e *quiz.TextEngine) error {
	if !e.LocateToNextSection() {
		Prt.Println("此位置不支持进行此跳转!!")
	}
	return nil
}

func skipToHead(e *quiz.TextEngine) error {
	e.LocateTo(e.HeadText)
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
	Prt.Println(Encouragement[rand.Intn(6)])
	return nil
}

func praise(e *quiz.TextEngine) error {
	Prt.Println(Praise[rand.Intn(6)])
	return nil
}

func printTittle(e *quiz.TextEngine) error {
	entry := e.CurrentEntry()
	if entry.IsTest {
		Prt.Print(entry.Tittle)
		return nil
	}
	Prt.Println(entry.Tittle)
	return nil
}

func setPath(e *quiz.TextEngine) error {
	text := e.CurrentText.(*TextBlock)
	pathKey := "tittles"
	var path []string
	p := e.UserCache[pathKey]
	if p != nil {
		path = p.([]string)
		for i, s := range path {
			if s == text.prev.tittle {
				path = path[:i+1]
				break
			}
		}
	}
	path = append(path, text.tittle)
	e.UserCache[pathKey] = path
	Prt.Prefix = strings.Join(path, "> ")
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
	fmt.Println(Prefix + "下面进入纠错模式")
	key := "wrongentrise"
	wrongEntrise := util.Load[map[*quiz.EntryQuiz]string](e.UserCache, key)
	if wrongEntrise == nil {
		return nil
	}
	entrise := util.Keys[*quiz.EntryQuiz](wrongEntrise)
	e.SetQuizEntrys(entrise)
	//在展示单个词条前设置打印前缀
	setPrefix := func(e *quiz.TextEngine) error {
		entry := e.CurrentEntry()
		Prt.Prefix = wrongEntrise[entry]
		return e.CheckEntry()
	}
	e.RegisterEntryFilter(setPrefix)
	e.ShowQuizEntrys()
	return nil
}
