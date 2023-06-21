package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
)

const (
	SkipKey  = "K"
	ScoreKey = "S"
	Quitkey  = "Q"
)

const (
	Welcome           = "欢迎使用子匀问答!\n计分请输入" + ScoreKey + "\n退出请输入" + Quitkey + "\n跳过请输入" + SkipKey
	GoodBye           = "就这样算了?不要怂啊!期待下次更好的你!"
	BeginQuiz         = "现在开始测验"
	TagQuestion       = `"%s"的谜底是：`
	AttentionQuestion = `"%s"的第%d个口诀是：`
	AnswerAgain       = "回答错误，请重新作答"
	NoAnswerBlank     = "作答不能为空!"
	PromptQuestion    = "请问是否需要提示"
	TimeOut           = "时间结束，你的分数是：%分"
	TransferIn        = "下面进入<%s>部分的测试!"
	TransferOut       = "恭喜你通过<%s>部分的测试!"
	Prefix            = "在%s中,"
)

var (
	Encouragement = [10]string{
		"再想想,你可以的!",
		"我不信你只有这种程度,再来!",
		"加油,坚持就是胜利!!!",
		"再做不出来就别吃饭了!",
		"这都错了啊,不是很简单的吗?",
		"答成这样可是要挨打的啊!",
	}
	Praise = [10]string{
		"不愧是你,这都做出来了!",
		"哈哈哈,不用扣鸡腿了!",
		"哟,不错嘛,又拿了一血!",
		"666!在下佩服!",
		"我不信你还能过关,再来!",
		"你这正确率,放在过去那就是状元!",
	}
)

var wd string

func startQuiz(tb *TextBlock) {
	fmt.Println(Welcome)
	fmt.Println(BeginQuiz)
	quizOld(tb)
}

func quizOld(tb *TextBlock) {
	EnterNewWd(tb)
	checkTag(tb)
	show(tb.statement)
	checkAttention(tb)
	show(tb.code)
	//time.Sleep(time.Second * 2)
	for i := 0; i < len(tb.subBlocks); i++ {
		sube := tb.subBlocks[i]
		quizOld(sube)
		ReturnOldWd(sube)
	}
}

func EnterNewWd(tb *TextBlock) bool {
	if len(tb.subBlocks) == 0 {
		return false
	}
	if wd == "" {
		wd = tb.name
		fmt.Println()
		fmt.Println(fmt.Sprintf(TransferIn, wd))
		return true
	}
	wd = wd + ">" + tb.name
	if tb.subBlocks[0].hasQuiz {
		fmt.Println()
		fmt.Println(fmt.Sprintf(TransferIn, strings.ReplaceAll(wd, ">", ".")))
	}
	return true
}

func ReturnOldWd(tb *TextBlock) bool {
	if len(tb.subBlocks) == 0 {
		return false
	}
	if tb.subBlocks[0].hasQuiz {
		fmt.Println()
		fmt.Println(fmt.Sprintf(TransferOut, strings.ReplaceAll(wd, ">", ".")))
	}
	wd = strings.TrimSuffix(wd, ">"+tb.name)
	return true
}

func show(s string) {
	if len(s) != 0 {
		fmt.Println()
		fmt.Println(s)
	}
}

func checkTag(tb *TextBlock) {
	if len(tb.tag) == 0 {
		return
	}
	question := strings.ReplaceAll(TagQuestion, "%s", tb.name)
	check(question, tb.tag)
}

func checkAttention(tb *TextBlock) {
	for i, v := range tb.attention {
		question := strings.ReplaceAll(AttentionQuestion, "%s", tb.name)
		question = fmt.Sprintf(question, i+1)
		check(question, v)
	}
}

func check(question, answer string) bool {
	var input string
	var prefix = fmt.Sprintf(Prefix, wd)
	for {
		fmt.Println()
		fmt.Print(prefix + question)
		n, err := fmt.Scan(&input)
		if err != nil {
			log.Println(err)
		}
		input = strings.TrimSpace(input)
		if n == 0 || len(input) == 0 {
			fmt.Println(NoAnswerBlank)
			return false
		}
		if HandleIfSkip(prefix+question, answer, input) {
			return false
		}
		HandleIfQuit(input)
		if input == answer {
			show(Praise[rand.Intn(6)])
			return true
		}
		show(Encouragement[rand.Intn(6)])
	}
}

func HandleIfSkip(question, answer, input string) bool {
	if input == SkipKey {
		fmt.Println()
		fmt.Print(question)
		fmt.Println(answer)
		return true
	}
	return false
}

func HandleIfQuit(s string) {
	if s == Quitkey {
		fmt.Println()
		fmt.Print(GoodBye)
		os.Exit(0)
	}
}
