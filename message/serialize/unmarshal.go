// SPDX-FileCopyrightText: 2020-2024 caixw
//
// SPDX-License-Identifier: MIT

package serialize

import (
	"io/fs"
	"os"
	"path/filepath"

	"github.com/issue9/localeutil"
	"github.com/issue9/sliceutil"

	"github.com/issue9/localeutil/message"
)

type UnmarshalFunc = func([]byte, any) error

// Search 根据文件名查找解码的方法
type Search = func(string) UnmarshalFunc

// Unmarshal 加载内容
func Unmarshal(data []byte, u UnmarshalFunc) (*message.Language, error) {
	l := &message.Language{}
	if err := u(data, l); err != nil {
		return nil, err
	}
	return l, nil
}

func LoadFile(path string, u UnmarshalFunc) (*message.Language, error) {
	return unmarshalFS(func() ([]byte, error) { return os.ReadFile(path) }, u)
}

func LoadFS(fsys fs.FS, name string, u UnmarshalFunc) (*message.Language, error) {
	return unmarshalFS(func() ([]byte, error) { return fs.ReadFile(fsys, name) }, u)
}

func unmarshalFS(f func() ([]byte, error), u UnmarshalFunc) (*message.Language, error) {
	data, err := f()
	if err != nil {
		return nil, err
	}
	return Unmarshal(data, u)
}

// LoadGlob 批量加载文件
//
// 相同语言 ID 的项会合并。
func LoadGlob(s Search, glob string) ([]*message.Language, error) {
	matches, err := filepath.Glob(glob)
	if err != nil {
		return nil, err
	}

	langs := make([]*message.Language, 0, len(matches))
	for _, match := range matches {
		u := s(match)
		if u == nil {
			return nil, localeutil.Error("not found unmarshal for %s", match)
		}
		l, err := LoadFile(match, u)
		if err != nil {
			return nil, err
		}
		langs = append(langs, l)
	}

	return joinLanguages(langs), nil
}

// LoadFSGlob 批量加载文件
//
// 相同语言 ID 的项会合并。
func LoadFSGlob(s Search, glob string, fsys ...fs.FS) ([]*message.Language, error) {
	langs := make([]*message.Language, 0, 10)
	for _, f := range fsys {
		matches, err := fs.Glob(f, glob)
		if err != nil {
			return nil, err
		}

		for _, match := range matches {
			u := s(match)
			if u == nil {
				return nil, localeutil.Error("not found unmarshal for %s", match)
			}

			l, err := LoadFS(f, match, u)
			if err != nil {
				return nil, err
			}
			langs = append(langs, l)
		}
	}

	return joinLanguages(langs), nil
}

func joinLanguages(langs []*message.Language) []*message.Language {
	delIndexes := make([]int, 0, len(langs))
	for index, lang := range langs {
		// 该元素已经被标记为删除
		if sliceutil.Exists(delIndexes, func(v int, _ int) bool { return index == v }) {
			continue
		}

		// 找与 lang.ID 相同的元素索引
		indexes := sliceutil.Indexes(langs, func(l *message.Language, i int) bool {
			return l.ID == lang.ID && i != index
		})

		for _, i := range indexes {
			lang.Join(langs[i])
		}

		delIndexes = append(delIndexes, indexes...)
	}

	return sliceutil.QuickDelete(langs, func(_ *message.Language, index int) bool {
		return sliceutil.Exists(delIndexes, func(i int, _ int) bool { return i == index })
	})
}
