// SPDX-FileCopyrightText: 2020-2024 caixw
//
// SPDX-License-Identifier: MIT

// Package localeutil 提供一些本地化的工具
package localeutil

import (
	"golang.org/x/text/language"

	"github.com/issue9/localeutil/internal/syslocale"
)

// DetectUserLanguage 检测当前用户的本地化信息
//
// 按以下顺序读取本地化信息：
//   - 环境变量 LANGUAGE；
//   - 平台相关，比如 windows 下调用 GetUserDefaultLocaleName 等；
//   - 按顺序读取 LC_ALL、LC_MESSAGES 和 LANG 环境变量；
//
// 所有的环境变量遵守 tag.encoding 的格式，比如 zh_CN.UTF-8，其中的 encoding 可以省略。
func DetectUserLanguage() string { return syslocale.Get() }

// DetectUserLanguageTag 检测当前用户的本地化信息
//
// 文档说明参考 [DetectUserLanguage]
func DetectUserLanguageTag() (language.Tag, error) {
	return language.Parse(syslocale.Get())
}
