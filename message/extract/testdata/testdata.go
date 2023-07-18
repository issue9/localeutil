// SPDX-License-Identifier: MIT

package testdata

import (
	"fmt"

	"github.com/issue9/localeutil"
)

var (
	_ = localeutil.Phrase("p1")
	_ = localeutil.Phrase("p1") // 同值，应该忽略
	_ = localeutil.Phrase("p1 %s", "str")
	_ = localeutil.Phrase("p1 %s:%d", "str", 5)
	_ = localeutil.Error("e1")
)

func f1() {
	fmt.Println(localeutil.Phrase("f1").LocaleString(nil))
	fmt.Println(localeutil.Phrase("f1 %d", 5).LocaleString(nil))

	// 变量，无法提取
	key := "k1 %d"
	arg := 5
	fmt.Println(localeutil.Phrase(key, arg).LocaleString(nil))
}

// 无法获取动态的内容
func f2(key string) {
	fmt.Println(localeutil.Phrase(key).LocaleString(nil))
}
