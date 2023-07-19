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
		LocaleString(*Printer) string
	}

	phrase struct {
		// NOTE: key 只能是字符串，如果要改为 [message.Reference]
		// 那么 message/extract 也要支持 [message.Key] 返回的所有类型。
		key    string
		values []any
	}

	localeError phrase

	// StringPhrase 由字符串组成的 [LocaleStringer] 实现
	//
	// 与 [Phrase] 不同，StringPhrase 可以是常量，且大部分情况下适用。
	StringPhrase string

	Printer = message.Printer
)

// Phrase 返回一段未翻译的语言片段
//
// key 和 val 参数与 [Printer.Sprintf] 的参数相同。
// 如果 val 也实现了 [LocaleStringer] 接口，则会先调用 val 的 LocaleString 方法。
//
// 如果 val 为空，将返回 StringPhrase(key)。
func Phrase(key string, val ...any) LocaleStringer {
	if len(val) == 0 {
		return StringPhrase(key)
	}
	return phrase{key: key, values: val}
}

// Error 返回未翻译的错误对象
//
// 该对象同时实现了 [LocaleStringer] 接口。
func Error(key string, val ...any) error {
	return &localeError{key: key, values: val}
}

func (p phrase) LocaleString(printer *Printer) string {
	values := make([]any, 0, len(p.values))
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

func (err *localeError) LocaleString(p *Printer) string {
	return phrase(*err).LocaleString(p)
}

func (sp StringPhrase) LocaleString(p *Printer) string { return p.Sprintf(string(sp)) }

func (sp StringPhrase) String() string { return sp.LocaleString(defaultPrinter) }
