// SPDX-License-Identifier: MIT

package localeutil

import "golang.org/x/text/message"

// LocaleStringer 本地化的字符串
type LocaleStringer interface {
	// LocaleString 返回当前对象的本地化字符串
	LocaleString(*message.Printer) string
}

// Phrase 返回一段未翻译的语言片段
func Phrase(key message.Reference, val ...interface{}) LocaleStringer {
	return phrase{key: key, values: val}
}

type phrase struct {
	key    message.Reference
	values []interface{}
}

func (p phrase) LocaleString(printer *message.Printer) string {
	return printer.Sprintf(p.key, p.values...)
}
