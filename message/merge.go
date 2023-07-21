// SPDX-License-Identifier: MIT

package message

import "github.com/issue9/sliceutil"

// Logger 日志输出接口
type Logger interface {
	Print(...any)
	Printf(string, ...any)
}

func mergeLanguage(src, dest *Language, log Logger) {
	if src.ID != dest.ID {
		return
	}

	// 删除只存在于 dest 而不存在于 src 的内容
	dest.Messages = sliceutil.Delete(dest.Messages, func(dm Message) bool {
		exist := sliceutil.Exists(src.Messages, func(sm Message) bool { return sm.Key == dm.Key })
		if !exist {
			log.Printf("删除语言 %s 的翻译项 %s", src.ID.String(), dm.Key)
		}
		return !exist
	})

	// 将 src 独有的项写入 dest
	for _, sm := range src.Messages {
		if !sliceutil.Exists(dest.Messages, func(dm Message) bool { return dm.Key == sm.Key }) {
			dest.Messages = append(dest.Messages, sm)
		}
	}
}

// Merge 将 src 并入当前对象
//
// 这将会执行以下几个步骤：
// - 删除只存在于 m 而不存在于 src 的内容；
// - 将 src 独有的项写入 dest；
// 最终内容是 dest 为准。
// log 所有删除的记录都将通过此输出；
func (m *Messages) Merge(src *Messages, log Logger) {
	// 删除只存在于 m 而不存在于 src 的内容
	m.Languages = sliceutil.Delete(m.Languages, func(dl *Language) bool {
		exist := sliceutil.Exists(src.Languages, func(sl *Language) bool { return sl.ID == dl.ID })
		if !exist {
			log.Printf("删除语言 %s", dl.ID.String())
		}
		return !exist
	})

	for index, ml := range m.Languages {
		sl, found := sliceutil.At(src.Languages, func(sl *Language) bool { return sl.ID == ml.ID })
		if !found {
			panic("在调用 Merge 时，外界修改了 Messages 内容")
		}
		mergeLanguage(sl, ml, log)
		m.Languages[index] = ml
	}

	// 将 src 独有的项写入 dest
	for _, sl := range src.Languages {
		if !sliceutil.Exists(m.Languages, func(dm *Language) bool { return dm.ID == sl.ID }) {
			m.Languages = append(m.Languages, sl)
		}
	}
}
