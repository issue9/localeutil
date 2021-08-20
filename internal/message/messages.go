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
		XMLName  struct{}        `xml:"messages" json:"-" yaml:"-"`
		Language language.Tag    `xml:"language,attr" json:"language" yaml:"language"`
		Messages []localeMessage `xml:"message" json:"messages" yaml:"messages"`
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

func LoadFromFS(b *catalog.Builder, fsys fs.FS, glob string, unmarshal func([]byte, interface{}) error) error {
	msgs, err := loadGlobFS(fsys, glob, unmarshal)
	if err != nil {
		return err
	}

	for _, msg := range msgs {
		if err := msg.set(b); err != nil {
			return err
		}
	}
	return nil
}

func loadGlobFS(fsys fs.FS, glob string, unmarshal func([]byte, interface{}) error) ([]*localeMessages, error) {
	matchs, err := fs.Glob(fsys, glob)
	if err != nil {
		return nil, err
	}

	msgs := make([]*localeMessages, 0, len(matchs))

	for _, file := range matchs {
		msg, err := loadMessages(fsys, file, unmarshal)
		if err != nil {
			return nil, err
		}
		msgs = append(msgs, msg)
	}

	return msgs, nil
}

func loadMessages(fsys fs.FS, file string, unmarshal func([]byte, interface{}) error) (*localeMessages, error) {
	data, err := fs.ReadFile(fsys, file)
	if err != nil {
		return nil, err
	}

	m := &localeMessages{}
	if err := unmarshal(data, m); err != nil {
		return nil, err
	}
	return m, nil
}

func (m *localeMessages) set(b *catalog.Builder) (err error) {
	tag := m.Language
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
