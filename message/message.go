// SPDX-License-Identifier: MIT

// Package message 本地化的语言文件处理
package message

import (
	"io/fs"
	"os"
	"path/filepath"

	"github.com/issue9/sliceutil"
	"golang.org/x/text/language"
)

type (
	// UnmarshalFunc 解析文本内容至对象的方法
	UnmarshalFunc = func([]byte, any) error

	MarshalFunc = func(any) ([]byte, error)

	// Messages 本地化对象
	Messages struct {
		XMLName   struct{}    `xml:"messages" json:"-" yaml:"-"`
		Languages []*Language `xml:"language" json:"languages" yaml:"languages"`
	}

	// Language 某一语言的本地化内容
	Language struct {
		ID       language.Tag `xml:"id,attr" json:"id" yaml:"id"`
		Messages []Message    `xml:"message" json:"messages" yaml:"messages"`
	}

	// Message 单条本地化内容
	Message struct {
		Key     string `xml:"key" json:"key" yaml:"key"`
		Message Text   `xml:"message" json:"message" yaml:"message"`
	}

	Text struct {
		Msg    string  `xml:"msg,omitempty" json:"msg,omitempty"  yaml:"msg,omitempty"`
		Select *Select `xml:"select,omitempty" json:"select,omitempty" yaml:"select,omitempty"`
		Vars   []*Var  `xml:"var,omitempty" json:"vars,omitempty" yaml:"vars,omitempty"`
	}

	Select struct {
		Arg    int     `xml:"arg,attr" json:"arg" yaml:"arg"`
		Format string  `xml:"format,attr,omitempty" json:"format,omitempty" yaml:"format,omitempty"`
		Cases  []*Case `xml:"case,omitempty" json:"cases,omitempty" yaml:"cases,omitempty"`
	}

	Var struct {
		Name   string  `xml:"name,attr" json:"name" yaml:"name"`
		Arg    int     `xml:"arg,attr" json:"arg" yaml:"arg"`
		Format string  `xml:"format,attr,omitempty" json:"format,omitempty" yaml:"format,omitempty"`
		Cases  []*Case `xml:"case,omitempty" json:"cases,omitempty" yaml:"cases,omitempty"`
	}

	Case struct {
		Case  string `xml:"case,attr" json:"case" yaml:"case"`
		Value string `xml:",chardata"`
	}
)

// Load 加载内容
func (m *Messages) Load(data []byte, u UnmarshalFunc) error {
	msgs := &Messages{}
	if err := u(data, msgs); err != nil {
		return err
	}
	m.merge(msgs)

	return nil
}

func (m *Messages) LoadFile(path string, u UnmarshalFunc) error {
	return m.unmarshalFS(func() ([]byte, error) { return os.ReadFile(path) }, u)
}

func (m *Messages) LoadFS(fsys fs.FS, name string, u UnmarshalFunc) error {
	return m.unmarshalFS(func() ([]byte, error) { return fs.ReadFile(fsys, name) }, u)
}

func (m *Messages) unmarshalFS(f func() ([]byte, error), u UnmarshalFunc) error {
	data, err := f()
	if err != nil {
		return err
	}
	return m.Load(data, u)
}

func (m *Messages) LoadGlob(glob string, u UnmarshalFunc) error {
	matches, err := filepath.Glob(glob)
	if err != nil {
		return err
	}

	for _, match := range matches {
		if err := m.LoadFile(match, u); err != nil {
			return err
		}
	}
	return nil
}

func (m *Messages) LoadFSGlob(fsys fs.FS, glob string, u UnmarshalFunc) error {
	matches, err := fs.Glob(fsys, glob)
	if err != nil {
		return err
	}

	for _, match := range matches {
		if err := m.LoadFS(fsys, match, u); err != nil {
			return err
		}
	}
	return nil
}

func (m *Messages) merge(m2 *Messages) {
	for _, l := range m2.Languages {
		ll, found := sliceutil.At(m.Languages, func(ll *Language) bool { return ll.ID == l.ID })
		if found {
			ll.merge(l)
		} else {
			m.Languages = append(m.Languages, l)
		}
	}
}

func (l *Language) merge(l2 *Language) {
	for _, msg2 := range l2.Messages {
		if !sliceutil.Exists(l.Messages, func(msg Message) bool { return msg.Key == msg2.Key }) {
			l.Messages = append(l.Messages, msg2)
		}
	}
}

// Bytes 将当前对象转换为 []byte
func (m *Messages) Bytes(f MarshalFunc) ([]byte, error) { return f(m) }

// SaveFile 将当前对象编码为文本并存入 path
func (m *Messages) SaveFile(dir, ext string, f MarshalFunc) error {
	if ext[0] != '.' {
		ext = "." + ext
	}

	msgs, err := m.marshal(f)
	if err != nil {
		return err
	}

	for k, v := range msgs {
		if err := os.WriteFile(filepath.Join(dir, k+ext), v, os.ModePerm); err != nil {
			return err
		}
	}

	return nil
}

func (m *Messages) marshal(f MarshalFunc) (map[string][]byte, error) {
	msgs := m.split()
	ret := make(map[string][]byte, len(msgs))

	for _, msg := range msgs {
		data, err := msg.Bytes(f)
		if err != nil {
			return nil, err
		}

		ret[msg.Languages[0].ID.String()] = data
	}

	return ret, nil
}

func (m *Messages) split() []*Messages {
	mm := make([]*Messages, 0, len(m.Languages))
	for _, l := range m.Languages {
		mm = append(mm, &Messages{Languages: []*Language{l}})
	}
	return mm
}
