// SPDX-License-Identifier: MIT

package localeutil

import (
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

var defaultPrinter = message.NewPrinter(language.Und)

type (
	// LocaleStringer 本地化字符串的接口中
	LocaleStringer interface {
		// LocaleString 返回当前对象的本地化字符串
		LocaleString(*message.Printer) string
	}

	phrase struct {
		key    message.Reference
		values []interface{}
	}

	localeError phrase
)

// Phrase 返回一段未翻译的语言片段
//
// key 和 val 参数与 golang.org/x/text/message.Printer.Sprintf 的参数相同。
// 如果 val 也实现了 LocaleStringer 接口，则会先调用 val 的 LocaleString 方法再传递给 Sprintf。
func Phrase(key message.Reference, val ...interface{}) LocaleStringer {
	return phrase{key: key, values: val}
}

// Error 返回未翻译的错误对象
//
// 该对象同时实现了 LocaleStringer 接口。
func Error(key message.Reference, val ...interface{}) error {
	return &localeError{key: key, values: val}
}

func (p phrase) LocaleString(printer *message.Printer) string {
	values := make([]interface{}, 0, len(p.values))
	for _, value := range p.values {
		if ls, ok := value.(LocaleStringer); ok {
			value = ls.LocaleString(printer)
		}
		values = append(values, value)
	}

	return printer.Sprintf(p.key, values...)
}

func (p phrase) String() string { return p.LocaleString(defaultPrinter) }

func (err *localeError) Error() string { return err.LocaleString(defaultPrinter) }

func (err *localeError) LocaleString(p *message.Printer) string {
	return phrase(*err).LocaleString(p)
}
