// Deprecated
package main

import (
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
)

const (
	CommonPrefixExp = "[-#*=` ]+"
	StatePrefixExp  = `\n[-#*= ]+[^-*\n]+-?\w*`
	HandleHeaderExp = `^[-#*= ]+(?P<name>[^-*\n]+)-?-?(?P<tag>[^\*\n]*)`
	AttentionExp    = `.*==(.+?)==`
	CodeExp         = "```" + `.*?(\n|\r)(?P<code>[\s|\S]+?)(\n|\r)` + "```"
)

var (
	CommonPrefixReg = regexp.MustCompile(CommonPrefixExp)
	StatePrefixReg  = regexp.MustCompile(StatePrefixExp)
	HeaderReg       = regexp.MustCompile(HandleHeaderExp)
	AttentionReg    = regexp.MustCompile(AttentionExp)
	CodeReg         = regexp.MustCompile(CodeExp)
)

// scan a file
func scan(s string) error {

	_, err := scanToelement(s)
	if err != nil {
		return err
	}
	return nil
}

func scanToelement(s string) (*TextBlock, error) {
	if Ext != path.Ext(s) {
		return nil, fmt.Errorf("ext of the fileName must be md")
	}
	content, err := os.ReadFile(s)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	_, fileName := filepath.Split(s)
	fileName = fileName[:strings.Index(fileName, ".")]
	e := TextBlock{name: fileName}
	//先格式化，去除冗余字符
	txt := clear(string(content))
	txt, code := getBackCode(txt)
	e.code = code
	splitByPrefix(txt, &e)
	return &e, nil
}

func splitByPrefix(txt string, e *TextBlock) {
	txt = strings.Trim(txt, "\n")
	prefix := getHeaderPrefix(txt)
	log.Println(prefix)
	if prefix == "" {
		log.Println("文本无法获取前缀")
		e.statement = txt
		return
	}
	reg2 := regexp.MustCompile(`(^|\s)` + regexp.QuoteMeta(prefix))
	indexes := reg2.FindAllStringIndex(txt, -1)
	subs := splitTxt(txt, indexes)
	for i := 0; i < len(subs); i++ {
		sube := handleTxt(subs[i])
		if sube != nil {
			sube.setHasQuiz()
			e.appendSub(sube)
		}
	}
}

func splitTxt(txt string, indexes [][]int) []string {
	l := len(indexes)
	subs := make([]string, len(indexes))
	start, end := indexes[0][0], -1
	for i := 1; i < l; i++ {
		end = indexes[i][0]
		subs[i-1] = txt[start:end]
		start = end
	}
	start = indexes[l-1][0]
	subs[l-1] = txt[start:]
	return subs
}

func handleTxt(txt string) (e *TextBlock) {
	txt = strings.Trim(txt, "\n\r")
	founded := HeaderReg.FindStringSubmatch(txt)
	if founded != nil {
		e = new(TextBlock)
		e.name = founded[HeaderReg.SubexpIndex("name")]
		e.tag = founded[HeaderReg.SubexpIndex("tag")]
	}
	//没有下一行就返回
	if e == nil || !strings.ContainsAny(txt, "\n\r") {
		return
	}
	//下移一行
	txt = txt[strings.IndexAny(txt, "\n\r")+1:]
	for {
		//头尾去除换行
		txt = strings.Trim(txt, "\n\r")
		if strings.TrimSpace(txt) == "" {
			return
		}
		//获取普通前缀
		prefix := getCommonPrefix(txt)
		switch {
		case strings.Index(prefix, CodeDelim) == 0:
			//取代码
			txt, e.code = getFrontCode(txt)

		case strings.Contains(prefix, "=="):
			//取要诀
			var a string
			txt, a = getAttention(txt)
			e.attention = append(e.attention, a)

		case !strings.ContainsAny(prefix, "*-#"):
			//取陈述
			txt, e.statement = getContent(txt)

		default:
			//将下一级分割
			splitByPrefix(txt, e)
			return e
		}
	}
}

func getContent(txt string) (string, string) {
	if !StatePrefixReg.MatchString(txt) {
		return "", txt
	}
	loc := StatePrefixReg.FindStringIndex(txt)
	return txt[loc[1]:], txt[:loc[0]]
}

func getAttention(txt string) (string, string) {
	var a string
	if AttentionReg.MatchString(txt) {
		subs := AttentionReg.FindStringSubmatchIndex(txt)
		a = txt[subs[2]:subs[3]]
		txt = txt[subs[1]:]
	}
	return txt, a
}

func getHeaderPrefix(txt string) string {
	return HeaderPrefixReg.FindString(txt)
}

func getCommonPrefix(txt string) string {
	return CommonPrefixReg.FindString(txt)
}

func getBackCode(txt string) (string, string) {
	var code string
	if strings.HasSuffix(txt, CodeDelim) {
		txt = strings.TrimSuffix(txt, CodeDelim)
		i := strings.LastIndex(txt, CodeDelim)
		code = txt[i+len(CodeDelim):]
		txt = txt[:i]
	}
	return txt, code
}

func getFrontCode(txt string) (string, string) {
	indexes := CodeReg.FindStringSubmatchIndex(txt)
	i := CodeReg.SubexpIndex("code") * 2
	start := indexes[i]
	end := indexes[i+1]
	return txt[indexes[1]:], txt[start:end]
}
