// SPDX-License-Identifier: MIT

package testdata

import (
	"fmt"

	"github.com/issue9/localeutil"
	"github.com/issue9/localeutil/testdata/locale"
)

const c1 = localeutil.StringPhrase("c1")

const constValue = "const-value"

var (
	_ = locale.String("c2")
	_ = localeutil.Phrase(constValue)

	_ = localeutil.Phrase("phrase 1")
	_ = localeutil.Phrase("phrase 1") // 同值，应该忽略
	_ = localeutil.Phrase("phrase %d", "2")

	_ = localeutil.Error("error 1")
	_ = localeutil.Error("error %d", 5)
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
