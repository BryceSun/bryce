package main

import (
	"example.com/bryce/quiz"
	"fmt"
	"log"
	"math/rand"
	"strings"
	"time"
)

const (
	SkipKey    = "K"
	ScoreKey   = "S"
	Quitkey    = "Q"
	PrefixFlag = "%s"
)

const (
	Welcome = "欢迎使用子匀问答!" +
		"\n计分请输入" + ScoreKey +
		"\n退出请输入" + Quitkey +
		"\n跳过请输入" + SkipKey
	GoodBye = "就这样算了???不要怂啊!\n期待下次更好的你!"
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
)

var engine *quiz.TextEngine

func showWith(doc *TextBlock) {
	engine = quiz.NewTextEngine(doc)
	engine.RegisterGuardFilter(showWelcome)
	engine.RegisterEntryFilter(showSpendedTime)
	engine.RegisterCoreFilter(setPath)
	engine.RegisterOrder("K", skip)
	engine.RegisterOrder("KK", skip1)
	engine.RegisterOrder("KKK", skip2)
	engine.RegisterOrder("Q", skipToHead)
	engine.RegisterOrder(quiz.EncourageFunKey, encourage)
	engine.RegisterOrder(quiz.PraiseFunKey, praise)
	engine.RegisterOrder(quiz.TittleFunKey, printTittle)
	engine.RegisterOrder("KH", skipToHead)
	err := engine.Start()
	if err != nil {
		log.Println(err)
	}
}
func showWelcome(e *quiz.TextEngine) error {
	fmt.Printf(PrefixFlag, Welcome+"\n")
	if err := e.Start(); err != nil {
		log.Println(err)
		return err
	}
	fmt.Printf(PrefixFlag, GoodBye+"\n")
	return nil
}

func skip(e *quiz.TextEngine) error {
	e.SetSkipOnce()
	s := getPath(e)
	entry := e.CurrentEntry()
	fmt.Printf(PrefixFlag, s+entry.Tittle+entry.Content+"\n")
	return nil
}

func skip1(e *quiz.TextEngine) error {
	e.LocateToNextText()
	return nil
}

func skip2(e *quiz.TextEngine) error {
	if !e.LocateToNextSection() {
		fmt.Println("此位置不支持进行此跳转!!")
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
		s := getPath(e)
		fmt.Printf(s+"本次作答用时:%d秒\n", int(time.Since(start).Seconds()))
		delete(e.UserCache, timeKey)
		return nil
	case e.HasLocate() || e.HasSkip():
		delete(e.UserCache, timeKey)
		return nil
	default:
		e.UserCache[timeKey] = start
		return nil
	}
}

func encourage(e *quiz.TextEngine) error {
	s := getPath(e)
	fmt.Printf(PrefixFlag, s+Encouragement[rand.Intn(6)]+"\n")
	return nil
}

func praise(e *quiz.TextEngine) error {
	s := getPath(e)
	fmt.Printf(PrefixFlag, s+Praise[rand.Intn(6)]+"\n")
	return nil
}

func printTittle(e *quiz.TextEngine) error {
	entry := e.CurrentEntry()
	s := getPath(e)
	if entry.IsTest {
		fmt.Printf(PrefixFlag, s+entry.Tittle)
		return nil
	}
	fmt.Printf(PrefixFlag, s+entry.Tittle+"\n")
	return nil
}

func getPath(e *quiz.TextEngine) string {
	p := e.UserCache["path"]
	if p != nil {
		return p.(string)
	}
	return ""
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
			}
		}
	}
	path = append(path, text.tittle)
	e.UserCache[pathKey] = path
	e.UserCache["path"] = strings.Join(path, "> ") + ">|  "
	return nil
}
