// SPDX-License-Identifier: MIT

package localeutil

import "golang.org/x/text/message"

type (
	// LocaleStringer 本地化字符串的接口中
	LocaleStringer interface {
		// LocaleString 返回当前对象的本地化字符串
		LocaleString(p *message.Printer) string
	}

	phrase struct {
		key    message.Reference
		values []interface{}
	}

	localeError struct {
		LocaleStringer
		p *message.Printer
	}
)

// Phrase 返回一段未翻译的语言片段
//
// key 和 val 参数与 golang.org/x/text/message.Printer.Sprintf 的参数相同。
// 如果 val 也实现了 LocaleStringer 接口，则会先调用 val 的 LocaleString 方法再传递给 Sprintf。
func Phrase(key message.Reference, val ...interface{}) LocaleStringer {
	return phrase{key: key, values: val}
}

// Error 构建错误对象
//
// 返回对象同时实现了 LocaleStringer 接口，用于对内容进行翻译。
//
// p 用于指定默认的语言对象，当用户直接调用 error.Error 方法时采用此对象；
func Error(p *message.Printer, key message.Reference, val ...interface{}) error {
	if p == nil {
		panic("p 不能为空")
	}
	return &localeError{p: p, LocaleStringer: Phrase(key, val...)}
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

func (err *localeError) Error() string { return err.LocaleString(err.p) }

func (err *localeError) LocaleString(p *message.Printer) string {
	return err.LocaleStringer.LocaleString(p)
}
