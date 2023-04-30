// SPDX-License-Identifier: MIT

package message

import (
	"io/fs"

	"golang.org/x/text/message/catalog"
)

// Load 从 data 解析本地化数据至 b
func Load(b *catalog.Builder, data []byte, unmarshal UnmarshalFunc) error {
	m := &Messages{}
	if err := unmarshal(data, m); err != nil {
		return err
	}
	return m.set(b)
}

// LoadFromFS 加载文件内容并写入 b
func LoadFromFS(b *catalog.Builder, fsys fs.FS, file string, unmarshal UnmarshalFunc) error {
	data, err := fs.ReadFile(fsys, file)
	if err != nil {
		return err
	}
	return Load(b, data, unmarshal)
}

// LoadFromFSGlob 加载多个文件内容并写入 b
func LoadFromFSGlob(b *catalog.Builder, fsys fs.FS, glob string, unmarshal UnmarshalFunc) error {
	matches, err := fs.Glob(fsys, glob)
	if err != nil {
		return err
	}

	for _, match := range matches {
		if err := LoadFromFS(b, fsys, match, unmarshal); err != nil {
			return err
		}
	}
	return nil
}
