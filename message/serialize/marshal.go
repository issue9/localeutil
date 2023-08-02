// SPDX-License-Identifier: MIT

package serialize

import (
	"io/fs"
	"os"
	"path/filepath"
	"sort"

	"github.com/issue9/localeutil/message"
)

type MarshalFunc = func(any) ([]byte, error)

// Marshal 将当前对象转换为 []byte
func Marshal(l *message.Language, f MarshalFunc) ([]byte, error) {
	// 输出前排序，保证相同内容输出的内容是一样的。
	sort.SliceStable(l.Messages, func(i, j int) bool {
		return l.Messages[i].Key < l.Messages[j].Key
	})

	return f(l)
}

// SaveFile 将当前对象编码为文本并存入 path
//
// 如果文件已经存在会被覆盖。
func SaveFile(l *message.Language, path string, f MarshalFunc, mode fs.FileMode) error {
	data, err := Marshal(l, f)
	if err == nil {
		err = os.WriteFile(path, data, mode)
	}
	return err
}

// SaveFiles 将当前对象按语言 ID 分类保存
func SaveFiles(langs []*message.Language, dir, ext string, f MarshalFunc, mode fs.FileMode) error {
	if ext[0] != '.' {
		ext = "." + ext
	}

	for _, l := range langs {
		path := filepath.Join(dir, l.ID.String()+ext)
		if err := SaveFile(l, path, f, mode); err != nil {
			return err
		}
	}
	return nil
}
