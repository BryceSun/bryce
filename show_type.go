package main

import (
	"encoding/json"
	"example.com/bryce/util"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"time"
)

var history map[string]string = make(map[string]string)
var folder string
var needSave bool

func LoadHistory(p string) error {
	folder = p
	f, err := os.OpenFile(folder+"\\type.log", os.O_CREATE|os.O_RDWR, os.ModePerm|os.ModeTemporary)
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			println(err)
			os.Exit(1)
		}
	}(f)
	if err != nil {
		panic(err)
	}
	h, err := io.ReadAll(f)
	if err != nil || len(h) == 0 {
		return err
	}
	return json.Unmarshal(h, &history)
}

func LogHistory() {
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Kill)
	ticker := time.Tick(time.Minute)
	for {
		select {
		case <-ticker:
			SaveHistory()
		case <-sigint:
			SaveHistory()
			os.Exit(0)
		}
	}
}

func SaveHistory() error {
	if !needSave {
		return nil
	}
	j, err := json.MarshalIndent(&history, "", "	")
	if err != nil {
		return err
	}
	err = os.WriteFile(folder+"\\type.log", j, 0666)
	needSave = false
	return err
}

func cachePosition(f, p string) {
	history[f] = p
	needSave = true
}

func clearPosition(f string) {
	history[f] = ""
	needSave = true
}

func getPosition(f string) string {
	return history[f]
}

func ShowType(filePath string) error {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}
	lines := strings.Split(string(content), "\n")
	target := ""
	jog := false
	position := ""
	locate := true
	if getPosition(filePath) == "" {
		locate = false
	}
LOOP:
	for _, line := range lines {
		l := strings.TrimSpace(line)
		if l == "" {
			continue
		}
		if locate {
			if l != getPosition(filePath) {
				continue
			}
			locate = false
		}
		if !jog || l == target || position == l {
			cachePosition(filePath, l)
			cmd := exec.Command("cmd", "/c", "cls")
			cmd.Stdout = os.Stdout
			err := cmd.Run()
			if err != nil {
				return err
			}
			jog = false
			fci := util.FirstChineseIndex(l)
			if fci > -1 {
				util.Rprintlnx(l[fci:])
				line = l[:fci]
			}
			util.Rprintln(line)
			util.Lprint(line)
			input := util.Scanln()
			if strings.Contains(input, "KK ") && len(input) > 3 {
				position = l
				target = strings.TrimPrefix(input, "KK ")
				jog = true
			}
		}
	}
	if target == "Q" {
		return nil
	}
	if jog == true {
		if target == "R" {
			target = ""
			jog = false
		}
		goto LOOP
	}
	clearPosition(filePath)
	return nil
}
