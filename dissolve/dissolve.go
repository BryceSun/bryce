package dissolve

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"example.com/bryce/util"
)

const IndentLinePattern = `(?:^|\n)(?P<%s>%s).*`

var Logger = log.Default()
var IndentRegs []*regexp.Regexp

type Tree interface {
	SetTittle(string)
	SetContent(string)
	NewTree() Tree
}

type textCube struct {
	Indent  string
	Tittle  string
	Content string
}

func (c textCube) fixTittle() {
	c.Tittle = strings.TrimPrefix(c.Tittle, c.Indent)

}

// ParseFile 根据文件名打开文件并解析生成树
func ParseFile(filePath string, t Tree) error {
	content, err := os.ReadFile(filePath)
	if err != nil {
		Logger.Println(err)
		return err
	}
	_, name := filepath.Split(filePath)
	name = name[:strings.Index(name, ".")]
	t.SetTittle(name)
	return SParse(string(content), t)
}

// FParse 直接从打开的文件解析并生成树
func FParse(file *os.File, t Tree) error {
	content, err := io.ReadAll(file)
	if err != nil {
		Logger.Println(err)
		return err
	}
	return SParse(string(content), t)
}

// SParse 直接从文本字符串中解析并生成树
func SParse(text string, t Tree) error {
	text = util.Clear(text)
	return transfer(text, t)
}

func transfer(text string, t Tree) error {
	content, textCubes := SplitText(text)
	if len(textCubes) == 0 {
		t.SetContent(content)
		return nil
	}
	err := transfer(content, t)
	if err != nil {
		return err
	}
	for _, cube := range textCubes {
		subt := t.NewTree()
		fixTextCube(cube)
		subt.SetTittle(cube.Tittle)
		err := transfer(cube.Content, subt)
		if err != nil {
			return err
		}
	}
	return nil
}

func fixTextCube(cube *textCube) {
	text, tittle := cube.Content, cube.Tittle
	cube.Tittle = strings.TrimPrefix(tittle, cube.Indent)
	Logger.Println("获取词条：", cube.Tittle)
	//没有下一行则返回
	lineEndIndex := strings.IndexAny(text, "\n\r")
	cube.Content = ""
	if lineEndIndex != -1 {
		cube.Content = text[lineEndIndex+1:]
	}
}

func SplitText(text string) (content string, subs []*textCube) {
	if strings.TrimSpace(text) == "" {
		return
	}
	var indentReg *regexp.Regexp
	for _, reg := range IndentRegs {
		if reg.MatchString(text) {
			indentReg = reg
			break
		}
	}
	if indentReg == nil {
		Logger.Println("文本不能被所有缩进正则表达式匹配")
		return text, subs
	}
	Logger.Println("使用正则表达式：", indentReg)
	firstindent := indentReg.FindString(text)
	patternName := "Indent"
	tittleLineReg := indentLineReg(firstindent, patternName)
	if !tittleLineReg.MatchString(text) {
		Logger.Panicln("文本不能被缩进正则表达式匹配, 表达式：", tittleLineReg)
		return
	}
	indexes := tittleLineReg.FindAllStringSubmatchIndex(text, -1)
	indentIndex := tittleLineReg.SubexpIndex(patternName) * 2

	l := len(indexes)
	subs = make([]*textCube, len(indexes))
	start, end := indexes[0][indentIndex], -1
	content = text[:start]
	for i := 0; i < l; i++ {
		//获取词条最大范围的文本
		j := i + 1
		subText := ""
		if j == l {
			subText = text[start:]
		} else {
			end = indexes[j][indentIndex]
			subText = text[start:end]
		}
		t := new(textCube)
		subs[i] = t
		t.Content = strings.Trim(subText, "\n\r")
		start = end
		//获取词条
		start2 := indexes[i][0]
		end2 := indexes[i][1]
		t.Tittle = strings.Trim(text[start2:end2], "\n\r")
		Logger.Printf("分割到第%d个词条%s", i, subs[i])
		//获取词条缩进
		start2 = indexes[i][indentIndex]
		end2 = indexes[i][indentIndex+1]
		t.Indent = strings.Trim(text[start2:end2], "\n\r")
	}
	return
}

func indentLineReg(indent, patternName string) *regexp.Regexp {
	pattern := fmt.Sprintf(IndentLinePattern, patternName, regexp.QuoteMeta(indent))
	indentReg, err := regexp.Compile(pattern)
	if err != nil {
		Logger.Panicln(err)
	}
	return indentReg
}
