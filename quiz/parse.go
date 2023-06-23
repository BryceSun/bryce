package quiz

import (
	"example.com/bryce/util"
	"fmt"
	"reflect"
	"strings"
)

const (
	// TagKey 字段标签中的键名
	TagKey = "quiz"
	VHide  = "hide"
	VCheck = "check"
	VHead  = "head"
	IFlag  = "${i}"
	KFlag  = "${k}"
)

type EntryQuiz struct {
	Tittle  string
	Content string
	IsTest  bool
}

func fieldMap(v reflect.Value) map[string]string {
	m := map[string]string{}
	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		ftype := t.Field(i)
		fvalue := v.Field(i)
		m[ftype.Name] = fvalue.String()
	}
	return m
}

// 解析QText字段生成词条
func parseQText(text QText) []*EntryQuiz {
	qtValue := reflect.ValueOf(text)
	return parse(qtValue)
}

// 截取标签标题
func getQuizTittle(tag string) string {
	i := strings.LastIndexByte(tag, '|')
	if i > 0 {
		return tag[i+1:]
	}
	return ""
}

// 解析结构体字段生成词条
func parse(v reflect.Value) []*EntryQuiz {
	var entries []*EntryQuiz
	if v.Kind() == reflect.Pointer {
		v = reflect.Indirect(v)
	}
	if v.Kind() != reflect.Struct {
		return nil
	}
	qtType := v.Type()
	xpMap := fieldMap(v)
	for i := 0; i < qtType.NumField(); i++ {
		field := qtType.Field(i)
		fieldValue := v.Field(i)
		tagv := field.Tag.Get(TagKey)
		if strings.TrimSpace(tagv) == "" || strings.Contains(tagv, VHide) || fieldValue.IsZero() {
			continue
		}
		tittle := getQuizTittle(tagv)
		tittle = util.Expand(tittle, xpMap)
		isTest := strings.Contains(tagv, VCheck)
		if strings.Contains(tagv, VHead) {
			entries = append(entries, &EntryQuiz{tittle, "", isTest})
			continue
		}
		entries = append(entries, getEntries(fieldValue, tittle, isTest)...)
	}
	return entries
}

// 根据变量值，词条标题，和测试标志生成词条
func getEntries(value reflect.Value, tittle string, isTest bool) []*EntryQuiz {
	var entries []*EntryQuiz
	tittle1 := tittle
	switch value.Kind() {
	case reflect.Struct:
		if strings.TrimSpace(tittle) != "" {
			entries = append(entries, &EntryQuiz{tittle, "", false})
		}
		return append(entries, parse(value)...)

	case reflect.Pointer:
		return getEntries(value.Elem(), tittle, isTest)

	case reflect.Map:
		r := value.MapRange()
		for r.Next() {
			k, v := r.Key().String(), r.Value().String()
			if strings.TrimSpace(tittle) != "" {
				tittle1 = strings.Replace(tittle, KFlag, k, 1)
			}
			entries = append(entries, &EntryQuiz{tittle1, v, isTest})
		}
		return entries

	case reflect.Array, reflect.Slice:
		for i := 0; i < value.Len(); i++ {
			v := value.Index(i)
			tittle1 = strings.Replace(tittle, IFlag, fmt.Sprint(i+1), 1)
			switch v.Kind() {
			case reflect.Struct:
				entries = append(entries, parse(value)...)
			case reflect.Pointer:
				entries = append(entries, getEntries(v.Elem(), tittle1, isTest)...)
			default:
				//todo
				entries = append(entries, getEntries(v, tittle1, isTest)...)
			}
		}
		return entries
	}
	return append(entries, &EntryQuiz{tittle, value.String(), isTest})
}
