// SPDX-License-Identifier: MIT

package testdata

import (
	"fmt"

	"github.com/issue9/localeutil"
	l "github.com/issue9/localeutil" // 为包指定别名
)

const c1 = localeutil.StringPhrase("c1")

const constValue = "const-value"

var (
	_ = localeutil.StringPhrase("c2")
	_ = l.Phrase(constValue)

	_ = l.Phrase("p1")
	_ = localeutil.Phrase("p1") // 同值，应该忽略
	_ = localeutil.Phrase("p1 %s", "str")

	_ = localeutil.Error("e1")
	_ = localeutil.Error("e1 %s:%d", "str", 5)
)

func f1() {
	fmt.Println(localeutil.Phrase("f1").LocaleString(nil))
	fmt.Println(localeutil.Phrase("f1 %d", 5).LocaleString(nil))

	// 变量，无法提取
	key := "k1 %d"
	arg := 5
	fmt.Println(localeutil.Phrase(key, arg).LocaleString(nil))

	// 需要计算的值，无法提取
	fmt.Println(localeutil.Phrase(constValue + "1").LocaleString(nil))
}

// 无法获取动态的内容
func f2(key string) {
	fmt.Println(localeutil.Phrase(key).LocaleString(nil))
}
