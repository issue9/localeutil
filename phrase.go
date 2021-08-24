// SPDX-License-Identifier: MIT

package localeutil

import "golang.org/x/text/message"

// Phrase 一段未翻译的语言片段
type Phrase struct {
	Key    message.Reference
	Values []interface{}
}

// LocaleString 返回当前语言的翻译内容
func (p *Phrase) LocaleString(printer *message.Printer) string {
	return printer.Sprintf(p.Key, p.Values...)
}
