// SPDX-License-Identifier: MIT

package localeutil

import (
	"fmt"

	"golang.org/x/text/message"
)

type (
	// LocaleStringer 本地化的字符串
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
func Phrase(key message.Reference, val ...interface{}) LocaleStringer {
	return phrase{key: key, values: val}
}

// Error 返回未翻译的错误对象
//
// 该对象同时实现了 LocaleStringer 接口。
func Error(key message.Reference, val ...interface{}) error {
	return localeError{key: key, values: val}
}

func (p phrase) LocaleString(printer *message.Printer) string {
	return printer.Sprintf(p.key, p.values...)
}

func (err localeError) Error() string { return fmt.Sprint(err.key) }

func (err localeError) LocaleString(p *message.Printer) string {
	return phrase(err).LocaleString(p)
}
