// SPDX-License-Identifier: MIT

// Package message 本地化的语言文件处理
package message

import (
	"io/fs"
	"os"

	"golang.org/x/text/language"
)

type (
	// UnmarshalFunc 解析文本内容至对象的方法
	UnmarshalFunc = func([]byte, interface{}) error

	MarshalFunc = func(interface{}) ([]byte, error)

	Messages struct {
		XMLName   struct{}       `xml:"messages" json:"-" yaml:"-"`
		Languages []language.Tag `xml:"language" json:"languages" yaml:"languages"`
		Messages  []Message      `xml:"message" json:"messages" yaml:"messages"`
	}

	// Message 单条消息
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
		Cases  []*Case `xml:"case" json:"cases" yaml:"cases"`
	}

	Var struct {
		Name   string  `xml:"name,attr" json:"name" yaml:"name"`
		Arg    int     `xml:"arg,attr" json:"arg" yaml:"arg"`
		Format string  `xml:"format,attr,omitempty" json:"format,omitempty" yaml:"format,omitempty"`
		Cases  []*Case `xml:"case" json:"cases" yaml:"cases"`
	}

	Case struct {
		Case  string `xml:"case,attr" json:"case" yaml:"case"`
		Value string `xml:",chardata"`
	}
)

func ex(cases []*Case) []interface{} {
	data := make([]interface{}, 0, len(cases)*2)
	for _, c := range cases {
		data = append(data, c.Case, c.Value)
	}
	return data
}

func Unmarshal(data []byte, u UnmarshalFunc) (*Messages, error) {
	m := &Messages{}
	err := u(data, m)
	return m, err
}

func UnmarshalFile(file string, u UnmarshalFunc) (*Messages, error) {
	return unmarshalFS(func() ([]byte, error) { return os.ReadFile(file) }, u)
}

func UnmarshalFS(fsys fs.FS, name string, u UnmarshalFunc) (*Messages, error) {
	return unmarshalFS(func() ([]byte, error) { return fs.ReadFile(fsys, name) }, u)
}

func unmarshalFS(f func() ([]byte, error), u UnmarshalFunc) (*Messages, error) {
	data, err := f()
	if err != nil {
		return nil, err
	}
	return Unmarshal(data, u)
}

func Marshal(f MarshalFunc, m *Messages) ([]byte, error) { return f(m) }

func MarshalFile(f MarshalFunc, m *Messages, path string) error {
	data, err := Marshal(f, m)
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, os.ModePerm)
}
