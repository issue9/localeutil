// SPDX-License-Identifier: MIT

// Package localeutil 提供一些本地化的工具
package localeutil

import (
	"io/fs"
	"os"
	"path/filepath"

	"golang.org/x/text/language"
	"golang.org/x/text/message/catalog"

	"github.com/issue9/localeutil/internal/message"
	"github.com/issue9/localeutil/internal/syslocale"
)

// DetectUserLanguageTag 检测当前用户的本地化信息
//
// *nix 系统会使用 LANG 环境变量中的值，windows 在 LANG
// 环境变量不存在的情况下，调用 GetUserDefaultLocaleName 函数获取。
func DetectUserLanguageTag() (language.Tag, error) { return syslocale.Get() }

// LoadMessageFromFS 从文件系统中加载文件并写入 b
//
// unmarshal 用于解析从 path 加载的文件；
func LoadMessageFromFS(b *catalog.Builder, fsys fs.FS, path string, unmarshal func([]byte, interface{}) error) error {
	return message.LoadFromFS(b, fsys, path, unmarshal)
}

// LoadMessageFromFile 从文件中加载文件并写入 b
func LoadMessageFromFile(b *catalog.Builder, path string, unmarshal func([]byte, interface{}) error) error {
	dir := filepath.ToSlash(filepath.Dir(path))
	path = filepath.ToSlash(filepath.Base(path))
	return LoadMessageFromFS(b, os.DirFS(dir), path, unmarshal)
}
