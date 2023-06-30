package util

import "fmt"

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
