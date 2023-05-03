// SPDX-License-Identifier: MIT

package message

import (
	"io/fs"

	"golang.org/x/text/feature/plural"
	"golang.org/x/text/message/catalog"
)

func (m *Messages) set(b *catalog.Builder) (err error) {
	for _, tag := range m.Languages {
		for _, msg := range m.Messages {
			switch {
			case msg.Message.Vars != nil:
				vars := msg.Message.Vars
				msgs := make([]catalog.Message, 0, len(vars))
				for _, v := range vars {
					mm := catalog.Var(v.Name, plural.Selectf(v.Arg, v.Format, ex(v.Cases)...))
					msgs = append(msgs, mm)
				}
				msgs = append(msgs, catalog.String(msg.Message.Msg))
				err = b.Set(tag, msg.Key, msgs...)
			case msg.Message.Select != nil:
				s := msg.Message.Select
				err = b.Set(tag, msg.Key, plural.Selectf(s.Arg, s.Format, ex(s.Cases)...))
			case msg.Message.Msg != "":
				err = b.SetString(tag, msg.Key, msg.Message.Msg)
			}

			if err != nil {
				return err
			}
		}
	}

	return nil
}

// Load 从 data 解析本地化数据至 b
func Load(b *catalog.Builder, data []byte, unmarshal UnmarshalFunc) error {
	m, err := Unmarshal(data, unmarshal)
	if err != nil {
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
