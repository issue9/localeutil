// SPDX-License-Identifier: MIT

package localeutil

import (
	"fmt"

	"golang.org/x/text/message"
)

type (
	// Stringer 本地化字符串
	Stringer interface {
		// LocaleString 返回当前对象的本地化字符串
		//
		// 如果 p 为 nil，可采用 [fmt.Sprintf] 的返回值。
		LocaleString(p *Printer) string
	}

	phrase struct {
		// NOTE: key 只能是字符串，如果要改为 [message.Reference]
		// 那么 message/extract 也要支持 [message.Key] 返回的所有类型。
		key    string
		values []any
	}

	localeError phrase

	// StringPhrase 由字符串组成的 [Stringer] 实现
	//
	// 与 [Phrase] 不同，StringPhrase 可以是常量，且大部分情况下适用。
	StringPhrase string

	Printer = message.Printer
)

// Phrase 返回一段未翻译的语言片段
//
// key 和 val 参数与 [Printer.Sprintf] 的参数相同。
// 如果 val 也实现了 [Stringer] 接口，则会先调用 val 的 LocaleString 方法。
//
// 如果 val 为空，将返回 StringPhrase(key)。
func Phrase(key string, val ...any) Stringer {
	if len(val) == 0 {
		return StringPhrase(key)
	}
	return phrase{key: key, values: val}
}

// Error 返回未翻译的错误对象
//
// 该对象同时实现了 [Stringer] 接口。
func Error(key string, val ...any) error {
	return &localeError{key: key, values: val}
}

func (p phrase) LocaleString(printer *Printer) string {
	if printer == nil {
		return fmt.Sprintf(p.key, p.values...)
	}

	values := make([]any, 0, len(p.values))
	for _, value := range p.values {
		if ls, ok := value.(Stringer); ok {
			value = ls.LocaleString(printer)
		}
		values = append(values, value)
	}

	return printer.Sprintf(p.key, values...)
}

func (err *localeError) Error() string { return phrase(*err).LocaleString(nil) }

func (err *localeError) LocaleString(p *Printer) string {
	return phrase(*err).LocaleString(p)
}

func (sp StringPhrase) LocaleString(p *Printer) string {
	if p == nil {
		return string(sp)
	}
	return p.Sprintf(string(sp))
}

// ErrorAsLocaleString 尝试将 err 转换为 [Stringer] 类型并输出
//
// 如果 err 未实现 [Stringer] 接口，则将调用 [error.Error]。
//
// NOTE: 未考虑 err 是否实现 xerrors.FormatError 等情况，
// 可作为简单的内容输出，正式环境作更多的类型判断。
func ErrorAsLocaleString(err error, p *Printer) string {
	if ls, ok := err.(Stringer); ok {
		return ls.LocaleString(p)
	}
	return err.Error()
}
