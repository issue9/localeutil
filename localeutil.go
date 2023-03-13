// SPDX-License-Identifier: MIT

// Package localeutil 提供一些本地化的工具
package localeutil

import (
	"io/fs"

	"golang.org/x/text/language"
	"golang.org/x/text/message/catalog"

	"github.com/issue9/localeutil/internal/message"
	"github.com/issue9/localeutil/internal/syslocale"
)

// UnmarshalFunc 解析文本内容至对象的方法
type UnmarshalFunc = message.UnmarshalFunc

// DetectUserLanguage 检测当前用户的本地化信息
//
// 默认会按顺序读取 LC_ALL、LC_MESSAGES 和 LANG 环境变量作为当前的语言环境。
// windows 在 LANG 环境变量不存在的情况下，调用 GetUserDefaultLocaleName 函数获取；
// darwin 会在 LANG 不存在的情况下，尝试读取 defaults read -g AppleLocale 中的值；
// js 平台则无视 LANG 环境变量，从 window.navigator.language 读取值；
func DetectUserLanguage() (string, error) { return syslocale.Get() }

// DetectUserLanguageTag 检测当前用户的本地化信息
//
// 文档说明参考 [DetectUserLanguage]
func DetectUserLanguageTag() (language.Tag, error) {
	l, err := syslocale.Get()
	if err != nil {
		return language.Und, err
	}
	return language.Parse(l)
}

// LoadMessage 解析 data 并写入 b
func LoadMessage(b *catalog.Builder, data []byte, unmarshal UnmarshalFunc) error {
	return message.Load(b, data, unmarshal)
}

// LoadMessageFromFS 从文件系统中加载文件并写入 b
//
// unmarshal 用于解析从 path 加载的文件；
func LoadMessageFromFS(b *catalog.Builder, fsys fs.FS, path string, unmarshal UnmarshalFunc) error {
	return message.LoadFromFS(b, fsys, path, unmarshal)
}

// LoadMessageFromFSGlob 从文件系统中加载多个文件并写入 b
func LoadMessageFromFSGlob(b *catalog.Builder, fsys fs.FS, glob string, unmarshal UnmarshalFunc) error {
	return message.LoadFromFSGlob(b, fsys, glob, unmarshal)
}
