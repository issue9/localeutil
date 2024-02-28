// SPDX-FileCopyrightText: 2020-2024 caixw
//
// SPDX-License-Identifier: MIT

package testdata

import (
	"fmt"

	"github.com/issue9/localeutil"
	"github.com/issue9/localeutil/testdata/locale"
	"github.com/issue9/localeutil/testdata/ref"
)

type String = localeutil.StringPhrase

const c1 = localeutil.StringPhrase("c1")

const c3 = String("c3")

const c4 = ref.String("c4")

const constValue = "const-value"

var p1 = &ref.Printer{}

var (
	_ = locale.String("c2")
	_ = localeutil.Phrase(constValue)

	_ = localeutil.Phrase("phrase 1")
	_ = localeutil.Phrase("phrase 1") // 同值，应该忽略
	_ = localeutil.Phrase("phrase %d", "2")

	_ = localeutil.Error("error 1")
	_ = localeutil.Error("error %d", 5)

	_ = p1.Print("testdata.Print")
)

func f1() {
	fmt.Println(localeutil.StringPhrase("f1").LocaleString(nil))
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
