// SPDX-License-Identifier: MIT

// Package localeutil 提供一些本地化的工具
package localeutil

import (
	"io/fs"
	"os"
	"path/filepath"

	"golang.org/x/text/language"
	"golang.org/x/text/message/catalog"

	"github.com/issue9/localeutil/internal/syslocale"
	"github.com/issue9/localeutil/message"
)

// SystemLanguageTag 返回当前系统的本地化信息
//
// *nix 系统会使用 LANG 环境变量中的值，windows 在 LANG
// 环境变量不存在的情况下，调用 GetUserDefaultLocaleName 函数获取。
func SystemLanguageTag() (language.Tag, error) {
	return syslocale.Get()
}

// LoadMessageFromFS 从文件系统中加载文件并写入 b
func LoadMessageFromFS(b *catalog.Builder, fsys fs.FS, glob string, unmarshal func([]byte, interface{}) error) error {
	return message.LoadFromFS(b, fsys, glob, unmarshal)
}

// LoadMessageFromFS 从文件中加载文件并写入 b
func LoadMessageFromFile(b *catalog.Builder, glob string, unmarshal func([]byte, interface{}) error) error {
	dir := filepath.ToSlash(filepath.Dir(glob))
	glob = filepath.ToSlash(filepath.Base(glob))
	return LoadMessageFromFS(b, os.DirFS(dir), glob, unmarshal)
}
