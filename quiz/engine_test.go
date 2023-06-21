package quiz

import (
	"fmt"
	"testing"
)

type qTextMocker struct {
	name      string   `quiz:"${name}的${tail}是:,head"`
	code      string   `quiz:"代码,false,hide"`
	tag       string   `quiz:"要诀,true,hide"`
	key       string   `quiz:"重点,true,show"`
	statement string   `quiz:"内容,false,show"`
	attention []string `quiz:"第${i}个口诀,true,show"`
	prev      *qTextMocker
	subs      []*qTextMocker
}

func (q *qTextMocker) Prev() QText {
	return q.prev
}

func (q *qTextMocker) Subs() []QText {
	var qts []QText
	for _, sub := range q.subs {
		qts = append(qts, sub)
	}
	return qts
}

var (
	attention0 = []string{}
	attention1 = []string{"one"}
	attention2 = []string{"one", "two"}
	attention3 = []string{"one", "two", "three"}
	attention4 = []string{"one", "two", "three", "for"}
	attention5 = []string{"one", "two", "three", "fou"}
	attention6 = []string{"one", "two", "three", "for"}
	atmh       = &qTextMocker{name: "nameh", subs: []*qTextMocker{atm0}}
	atm0       = &qTextMocker{"name0", "hello,world0", "hello0", "key1", "statement0", attention0, nil, []*qTextMocker{atm1, atm2}}
	atm1       = &qTextMocker{"name1", "hello,world1", "hello1", "key2", "statemen1", attention1, nil, []*qTextMocker{atm3, atm4}}
	atm2       = &qTextMocker{"name2", "hello,world2", "hello2", "key3", "statement2", attention2, nil, []*qTextMocker{atm5, atm6}}
	atm3       = &qTextMocker{"name3", "hello,world3", "hello3", "key4", "statement3", attention3, nil, nil}
	atm4       = &qTextMocker{"name4", "hello,world4", "hello4", "key5", "statement4", attention4, nil, nil}
	atm5       = &qTextMocker{"name5", "hello,world5", "hello5", "key6", "statement5", attention5, nil, nil}
	atm6       = &qTextMocker{"name6", "hello,world6", "hello6", "key6", "statement6", attention6, nil, nil}
)

func initAtm() {
	atm0.prev = atmh
	atm1.prev = atm0
	atm2.prev = atm0
	atm3.prev = atm1
	atm4.prev = atm1
	atm5.prev = atm2
	atm6.prev = atm2
}

func TestPrseQText(t *testing.T) {
	entryQuizs := parseQText(atm1)
	fmt.Printf("%+v", entryQuizs)
}

func TestNewTextEngine(t *testing.T) {
	engine := NewTextEngine(atmh)
	fmt.Printf("%+v", engine)
	showWelcome(engine)
	showGoodBye(engine)
	TestTextEngine_RegisterHandler(t)
}

func TestTextEngine_RegisterHandler(t *testing.T) {
	engine := NewTextEngine(atmh)
	userOrders := []Handler{func(textEngine *TextEngine) error {
		fmt.Println("用户命令1。。。。")
		return nil
	},
		func(textEngine *TextEngine) error {
			fmt.Println("用户命令2。。。。")
			return nil
		}}
	defaultOrders := []Handler{func(textEngine *TextEngine) error {
		fmt.Println("默认命令1。。。。")
		return nil
	}, func(textEngine *TextEngine) error {
		fmt.Println("默认命令2。。。。")
		return nil
	}}
	one, two := "order1", "order2"
	_ = userOrders
	_ = defaultOrders
	engine.RegisterHandler(one, userOrders[0])
	engine.RegisterHandler(two, userOrders[1])
	engine.registerHandler(one, defaultOrders[0])
	engine.registerHandler(two, defaultOrders[1])
	engine.excFuncOrPrintln(one, "没有注册"+one+"功能")
	engine.excFuncOrPrintln(two, "没有注册"+two+"功能")
	engine.excFuncOrPrintln("onefunc", "没有注册onefunc功能")
	initAtm()
	engine.SetToText(atm1)
	if engine.SetToRighText() != true || engine.currentText != atm2 {
		t.Fail()
	}
	if engine.SetToRighText() != false {
		t.Fail()
	}
	if engine.SetToUpperText() != true || engine.currentText != atm0 {
		t.Fail()
	}
	if engine.SetToRighText() != false {
		t.Fail()
	}
	if engine.SetToUpperText() != true || engine.currentText != atmh {
		t.Fail()
	}
	if engine.SetToRighText() != false || engine.SetToUpperText() != false || engine.currentText != atmh {
		t.Fail()
	}
	println(engine)
}

func TestTextEngine_Start(t *testing.T) {
	engine := NewTextEngine(atmh)
	engine.Start()
}
