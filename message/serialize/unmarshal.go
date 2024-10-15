// SPDX-FileCopyrightText: 2020-2024 caixw
//
// SPDX-License-Identifier: MIT

package serialize

import (
	"io/fs"
	"os"
	"path/filepath"

	"github.com/issue9/localeutil"
	"github.com/issue9/localeutil/message"
)

type UnmarshalFunc = func([]byte, any) error

// Search 根据文件名查找解码的方法
type Search = func(string) UnmarshalFunc

// Unmarshal 加载内容
func Unmarshal(data []byte, u UnmarshalFunc) (*message.File, error) {
	l := &message.File{}
	if err := u(data, l); err != nil {
		return nil, err
	}
	return l, nil
}

func LoadFile(path string, u UnmarshalFunc) (*message.File, error) {
	return unmarshalFS(func() ([]byte, error) { return os.ReadFile(path) }, u)
}

func LoadFS(fsys fs.FS, name string, u UnmarshalFunc) (*message.File, error) {
	return unmarshalFS(func() ([]byte, error) { return fs.ReadFile(fsys, name) }, u)
}

func unmarshalFS(f func() ([]byte, error), u UnmarshalFunc) (*message.File, error) {
	data, err := f()
	if err != nil {
		return nil, err
	}
	return Unmarshal(data, u)
}

// LoadGlob 批量加载文件
func LoadGlob(s Search, glob string) ([]*message.File, error) {
	matches, err := filepath.Glob(glob)
	if err != nil {
		return nil, err
	}

	langs := make([]*message.File, 0, len(matches))
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

	return langs, nil
}

// LoadFSGlob 批量加载文件
func LoadFSGlob(s Search, glob string, fsys ...fs.FS) ([]*message.File, error) {
	langs := make([]*message.File, 0, 10)
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

	return langs, nil
}
