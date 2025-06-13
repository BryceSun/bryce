package util

import (
	"fmt"
	"github.com/buger/goterm"
	"strconv"
	"unicode/utf8"
)

type Printer struct {
	Prefix string
	Format string
}

var printer = Printer{Format: "%s%s"}

func NewPrinter(prefix string, format string) *Printer {
	return &Printer{Prefix: prefix, Format: format}
}

func (p *Printer) Printf(format string, a ...any) {
	s := fmt.Sprintf(format, a...)
	fmt.Printf(p.Format, p.Prefix, s)
}

func (p *Printer) Println(a ...any) {
	s := fmt.Sprintln(a...)
	fmt.Printf(p.Format, p.Prefix, s)
}

func (p *Printer) Print(a ...any) {
	s := fmt.Sprint(a...)
	fmt.Printf(p.Format, p.Prefix, s)
}

func (p *Printer) Sprint(a ...any) string {
	s := fmt.Sprint(a...)
	return fmt.Sprintf(p.Format, p.Prefix, s)
}

func Printf(format string, a ...any) {
	printer.Printf(format, a...)
}

func Println(a ...any) {
	printer.Println(a...)
}

func Print(a ...any) {
	printer.Print(a...)
}

func Sprint(a ...any) string {
	return printer.Sprint(a...)
}

func Rprintln(s string) {
	terminalWidth := goterm.Width()
	n := utf8.RuneCountInString(s)
	posit := strconv.Itoa((terminalWidth + n) / 2)
	fmt.Printf("%"+posit+"s\n", s)
}

func Rprintlnx(s string) {
	terminalWidth := goterm.Width()
	n := utf8.RuneCountInString(s)
	posit := strconv.Itoa((terminalWidth+n/2)/2 - 2)
	fmt.Printf("%"+posit+"s\n", s)
}

func Lprint(s string) {
	terminalWidth := goterm.Width()
	n := utf8.RuneCountInString(s)
	posit := strconv.Itoa((terminalWidth - n) / 2)
	fmt.Printf("%"+posit+"s", "")
}
