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

type UnmarshalFunc = func([]byte, interface{}) error

// DetectUserLanguageTag 检测当前用户的本地化信息
//
// 使用 LANG 环境变量中的值，windows 在 LANG
// 环境变量不存在的情况下，调用 GetUserDefaultLocaleName 函数获取。
func DetectUserLanguageTag() (language.Tag, error) { return syslocale.Get() }

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
