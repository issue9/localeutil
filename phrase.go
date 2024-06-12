// SPDX-FileCopyrightText: 2020-2024 caixw
//
// SPDX-License-Identifier: MIT

package localeutil

import (
	"errors"
	"fmt"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"golang.org/x/text/message/catalog"
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

	phraseError phrase

	stringError struct {
		// 不能是常量，否则无法处理 [errors.Is] 和 [errors.As] 等操作
		key string
	}

	// StringPhrase 由字符串组成的 [Stringer] 实现
	//
	// 与 [Phrase] 不同，StringPhrase 是常量，且大部分情况下适用。
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
// 如果 val 中包含 error 对象，可以用 [errors.Is] 进行检测。
func Error(key string, val ...any) error {
	if len(val) == 0 {
		return &stringError{key: key}
	}
	return &phraseError{key: key, values: val}
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

func (err *phraseError) Error() string { return err.LocaleString(nil) }

func (err *phraseError) LocaleString(p *Printer) string {
	return phrase(*err).LocaleString(p)
}

func (err *phraseError) Is(target error) bool {
	for _, v := range err.values {
		if e, ok := v.(error); ok && errors.Is(e, target) {
			return true
		}
	}
	return false
}

func (err *stringError) Error() string { return err.LocaleString(nil) }

func (err *stringError) LocaleString(p *Printer) string {
	return StringPhrase(err.key).LocaleString(p)
}

func (sp StringPhrase) LocaleString(p *Printer) string {
	if p == nil {
		return string(sp)
	}
	return p.Sprintf(string(sp))
}

// NewPrinter 从 cat 查找最符合 tag 的语言作为打印对象返回
func NewPrinter(cat catalog.Catalog, tag language.Tag) *Printer {
	tag, _, _ = cat.Matcher().Match(tag)
	return message.NewPrinter(tag, message.Catalog(cat))
}
