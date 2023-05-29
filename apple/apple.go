// SPDX-License-Identifier: MIT

//go:build darwin || ios

// Package apple 苹果部分系统下的专有接口
package apple

import (
	"strings"
	"unicode"

	"github.com/issue9/localeutil/internal/defaults"
	"github.com/issue9/localeutil/internal/syslocale"
)

const appleLanguagesKey = "AppleLanguages"

// AppLocale 返回 app 的界面语言
//
// app 表示该应用的 ID；
//
// NOTE: macOS 系统中可以在设置中修改每个应用的语言，该接口可以获取此值。
func AppLocale(app string) string {
	v := defaults.Read(appleLanguagesKey, app)

	langs := strings.Split(strings.Trim(v, "()"), ",")
	if len(langs) == 0 {
		return syslocale.Get()
	}
	return strings.TrimFunc(langs[0], func(r rune) bool {
		return r == '"' || unicode.IsSpace(r)
	})
}

// SetAppLocale 设置 app 的界面语言
//
// app 表示应用有的唯一 ID；
// lang 为语言的 ID，如果是多个值，那么是第一个真实存在的 ID 作为其实际语言；
//
// NOTE: macOS 系统中可以在设置中修改每个应用的语言，该接口可以设置此值。
func SetAppLocale(app string, lang ...string) error {
	if len(lang) == 0 {
		panic("参数 lang 不能为空")
	}

	return defaults.Write(app, appleLanguagesKey, "-array", lang...)
}
