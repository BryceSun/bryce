package quiz

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
)

const (
	// TagKey 字段标签中的键名
	TagKey    = "quiz"
	VHead     = "head"
	VHide     = "hide"
	VTest     = "true"
	NameKey   = "tail"
	IndexIKey = "i"
	IndexJKey = "j"
)

type EntryQuiz struct {
	tittle  string
	content string
	isTest  bool
}

func parseQText(text QText) []*EntryQuiz {
	var entrys []*EntryQuiz
	testMap := map[string]string{}
	stateMap := map[string]string{}
	qtValue := reflect.ValueOf(text)
	if qtValue.Kind() == reflect.Ptr {
		qtValue = reflect.Indirect(qtValue)
	}
	if qtValue.Kind() != reflect.Struct {
		return nil
	}
	tbType := qtValue.Type()
	numField := tbType.NumField()
	quizHead := ""
	for i := 0; i < numField; i++ {
		field := tbType.Field(i)
		fieldValue := qtValue.Field(i)
		tag := field.Tag
		tagv := tag.Get(TagKey)
		if tagv == "" {
			continue
		}
		quizName := getQuizName(tagv)
		if strings.Contains(tagv, VHead) {
			quizHead = quizName
			loc := fmt.Sprintf("${%s}", field.Name)
			quizHead = strings.ReplaceAll(quizHead, loc, fieldValue.String())
			continue
		}
		if strings.Contains(tagv, VHide) {
			continue
		}
		if strings.Contains(tagv, VTest) {
			setQuizMap(testMap, quizHead, quizName, fieldValue)
		} else {
			setQuizMap(stateMap, quizHead, quizName, fieldValue)
		}
	}
	entrys = append(entrys, transferToEntrys(stateMap, false)...)
	entrys = append(entrys, transferToEntrys(testMap, true)...)
	return entrys
}

func transferToEntrys(quizMap map[string]string, isTest bool) []*EntryQuiz {
	var entrys []*EntryQuiz
	for k, v := range quizMap {
		if strings.TrimSpace(v) != "" {
			entrys = append(entrys, &EntryQuiz{k, v, isTest})
		}
	}
	return entrys
}

func getQuizName(tag string) string {
	i := strings.IndexRune(tag, ',')
	if i > 0 {
		return tag[:i]
	}
	return tag
}

func setQuizMap(quizMap map[string]string, quizHead string, quizName string, value reflect.Value) {
	i := 0
	expandfunc := func(s string) string {
		switch s {
		case NameKey:
			return quizName
		case IndexIKey, IndexJKey:
			return strconv.Itoa(i + 1)
		default:
			return ""
		}
	}
	setMap := func(name string, answer string) {
		quiz := os.Expand(name, expandfunc)
		quizMap[quiz] = answer
	}
	switch value.Kind() {
	case reflect.String:
		setMap(quizHead, value.String())
	case reflect.Array, reflect.Slice:
		for ; i < value.Len(); i++ {
			v := value.Index(i)
			quizHead := os.Expand(quizHead, expandfunc)
			if v.Kind() == reflect.String {
				setMap(quizHead, v.String())
			}
			if v.Kind() == reflect.Array {
				setQuizMap(quizMap, quizHead, quizName, v)
			}
		}
	}
}
