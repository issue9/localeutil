// SPDX-License-Identifier: MIT

// Package message 从文件中加载本地化信息
package message

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"errors"
	"io"
	"io/fs"

	"golang.org/x/text/feature/plural"
	"golang.org/x/text/language"
	"golang.org/x/text/message/catalog"
	"gopkg.in/yaml.v2"
)

type (
	localeMessages struct {
		XMLName   struct{}        `xml:"messages" json:"-" yaml:"-"`
		Languages []language.Tag  `xml:"language" json:"languages" yaml:"languages"`
		Messages  []localeMessage `xml:"message" json:"messages" yaml:"messages"`
	}

	// localeMessage 单条消息
	localeMessage struct {
		Key     string     `xml:"key" json:"key" yaml:"key"`
		Message localeText `xml:"message" json:"message" yaml:"message"`
	}

	localeText struct {
		Msg    string        `xml:"msg,omitempty" json:"msg,omitempty"  yaml:"msg,omitempty"`
		Select *localeSelect `xml:"select,omitempty" json:"select,omitempty" yaml:"select,omitempty"`
		Vars   []*localeVar  `xml:"var,omitempty" json:"vars,omitempty" yaml:"vars,omitempty"`
	}

	localeSelect struct {
		Arg    int         `xml:"arg,attr" json:"arg" yaml:"arg"`
		Format string      `xml:"format,attr,omitempty" json:"format,omitempty" yaml:"format,omitempty"`
		Cases  localeCases `xml:"case" json:"cases" yaml:"cases"`
	}

	localeVar struct {
		Name   string      `xml:"name,attr" json:"name" yaml:"name"`
		Arg    int         `xml:"arg,attr" json:"arg" yaml:"arg"`
		Format string      `xml:"format,attr,omitempty" json:"format,omitempty" yaml:"format,omitempty"`
		Cases  localeCases `xml:"case" json:"cases" yaml:"cases"`
	}

	localeCases []interface{}

	localeCaseEntry struct {
		XMLName struct{} `xml:"case"`
		Cond    string   `xml:"cond,attr"`
		Value   string   `xml:",chardata"`
	}
)

// LoadFromFS 加载文件内容并写入 b
func LoadFromFS(b *catalog.Builder, fsys fs.FS, file string, unmarshal func([]byte, interface{}) error) error {
	data, err := fs.ReadFile(fsys, file)
	if err != nil {
		return err
	}

	m := &localeMessages{}
	if err := unmarshal(data, m); err != nil {
		return err
	}

	return m.set(b)
}

// LoadFromFSGlob 加载多个文件内容并写入 b
func LoadFromFSGlob(b *catalog.Builder, fsys fs.FS, glob string, unmarshal func([]byte, interface{}) error) error {
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

func (m *localeMessages) set(b *catalog.Builder) (err error) {
	for _, tag := range m.Languages {
		for _, msg := range m.Messages {
			switch {
			case msg.Message.Vars != nil:
				vars := msg.Message.Vars
				msgs := make([]catalog.Message, 0, len(vars))
				for _, v := range vars {
					mm := catalog.Var(v.Name, plural.Selectf(v.Arg, v.Format, v.Cases...))
					msgs = append(msgs, mm)
				}
				msgs = append(msgs, catalog.String(msg.Message.Msg))
				err = b.Set(tag, msg.Key, msgs...)
			case msg.Message.Select != nil:
				s := msg.Message.Select
				err = b.Set(tag, msg.Key, plural.Selectf(s.Arg, s.Format, s.Cases...))
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

// UnmarshalXML implement xml.Unmarshaler
func (c *localeCases) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	for {
		e := &localeCaseEntry{}
		if err := d.DecodeElement(e, &start); errors.Is(err, io.EOF) {
			return nil
		} else if err != nil {
			return err
		}

		*c = append(*c, e.Cond, e.Value)
	}
}

func (c *localeCases) UnmarshalYAML(unmarshal func(interface{}) error) error {
	kv := yaml.MapSlice{}
	if err := unmarshal(&kv); err != nil {
		return err
	}

	*c = make(localeCases, 0, len(kv))
	for _, item := range kv {
		*c = append(*c, item.Key, item.Value)
	}

	return nil
}

func (c *localeCases) UnmarshalJSON(data []byte) error {
	d := json.NewDecoder(bytes.NewBuffer(data))
	for {
		t, err := d.Token()
		if errors.Is(err, io.EOF) {
			return nil
		} else if err != nil {
			return err
		}

		if t == json.Delim('{') || t == json.Delim('}') {
			continue
		}

		*c = append(*c, t)
	}
}
