package main

import (
	"example.com/bryce/quiz"
	"fmt"
	"math/rand"
	"time"
)

const (
	SkipKey  = "K"
	ScoreKey = "S"
	Quitkey  = "Q"
)

const (
	Welcome = "欢迎使用子匀问答!" +
		"\n计分请输入" + ScoreKey +
		"\n退出请输入" + Quitkey +
		"\n跳过请输入" + SkipKey + "\n"
	GoodBye = "就这样算了???不要怂啊!\n期待下次更好的你!\n"
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
	engine.RegisterOrder("K", skip)
	engine.RegisterOrder("KK", skip1)
	engine.RegisterOrder("KKK", skip2)
	engine.RegisterOrder("KH", skipToHead)
	engine.RegisterOrder(quiz.EncourageFunKey, encourage)
	engine.RegisterOrder(quiz.PraiseFunKey, praise)
	engine.RegisterOrder("KH", skipToHead)
	engine.Start()
}
func showWelcome(e *quiz.TextEngine) error {
	fmt.Printf(Welcome)
	e.Start()
	fmt.Printf(GoodBye)
	return nil
}

func skip(e *quiz.TextEngine) error {
	e.SetSkipOnce()
	entry := e.CurrentEntry()
	fmt.Println(entry.Tittle, entry.Content)
	return nil
}

func skip1(e *quiz.TextEngine) error {
	e.LocateToNextText()
	return nil
}

func skip2(e *quiz.TextEngine) error {
	e.LocateToNextSection()
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
	e.CheckEntry()
	switch {
	case e.Right:
		fmt.Printf("本次作答用时:%d秒\n", int(time.Since(start).Seconds()))
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
	fmt.Println(Encouragement[rand.Intn(6)])
	return nil
}

func praise(e *quiz.TextEngine) error {
	fmt.Println(Praise[rand.Intn(6)])
	return nil
}
