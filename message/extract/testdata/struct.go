// SPDX-License-Identifier: MIT

package testdata

import (
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	xm "golang.org/x/text/message"
)

var globalVar message.Printer

func structFunc(fp xm.Printer) {
	var local *xm.Printer = xm.NewPrinter(language.SimplifiedChinese)
	local.Printf("local var")

	globalVar.Printf("global var")
	fp.Printf("field var")

	// 无法完成类型推导
	s := message.NewPrinter(language.SimplifiedChinese)
	s.Printf("s1")
}
