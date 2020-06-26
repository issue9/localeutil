// SPDX-License-Identifier: MIT

// Package localeutil 提供一些本地化的工具
package localeutil

import (
	"golang.org/x/text/language"

	"github.com/issue9/localeutil/internal/syslocale"
)

// SystemLanguageTag 返回当前系统的本地化信息
//
// *nix 系统会使用 LANG 环境变量中的值，windows 在 LANG
// 环境变量不存在的情况下，调用 GetUserDefaultLocaleName 函数获取。
func SystemLanguageTag() (language.Tag, error) {
	return syslocale.Get()
}
