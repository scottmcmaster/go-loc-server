package loader

import (
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

func getPrinter(lang string) *message.Printer {
	tag, _ := language.Parse(lang)
	return message.NewPrinter(tag)
}
