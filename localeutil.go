// SPDX-License-Identifier: MIT

// Package localeutil 提供一些本地化的工具
package localeutil

import (
	"io/fs"

	"golang.org/x/text/language"
	"golang.org/x/text/message/catalog"

	"github.com/issue9/localeutil/internal/syslocale"
	"github.com/issue9/localeutil/message"
)

// UnmarshalFunc 解析文本内容至对象的方法
type UnmarshalFunc = message.UnmarshalFunc

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
