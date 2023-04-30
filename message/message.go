// SPDX-License-Identifier: MIT

// Package message 从文件中加载本地化信息
package message

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"errors"
	"io"

	"golang.org/x/text/feature/plural"
	"golang.org/x/text/language"
	"golang.org/x/text/message/catalog"
	"gopkg.in/yaml.v3"
)

type (
	// UnmarshalFunc 解析文本内容至对象的方法
	UnmarshalFunc = func([]byte, interface{}) error

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
		Arg    int    `xml:"arg,attr" json:"arg" yaml:"arg"`
		Format string `xml:"format,attr,omitempty" json:"format,omitempty" yaml:"format,omitempty"`
		Cases  Cases  `xml:"case" json:"cases" yaml:"cases"`
	}

	Var struct {
		Name   string `xml:"name,attr" json:"name" yaml:"name"`
		Arg    int    `xml:"arg,attr" json:"arg" yaml:"arg"`
		Format string `xml:"format,attr,omitempty" json:"format,omitempty" yaml:"format,omitempty"`
		Cases  Cases  `xml:"case" json:"cases" yaml:"cases"`
	}

	Cases []interface{}

	caseEntry struct {
		XMLName struct{} `xml:"case"`
		Cond    string   `xml:"cond,attr"`
		Value   string   `xml:",chardata"`
	}
)

func (m *Messages) set(b *catalog.Builder) (err error) {
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

func (c *Cases) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	for {
		e := &caseEntry{}
		if err := d.DecodeElement(e, &start); errors.Is(err, io.EOF) {
			return nil
		} else if err != nil {
			return err
		}

		*c = append(*c, e.Cond, e.Value)
	}
}

func (c *Cases) UnmarshalYAML(value *yaml.Node) error {
	l := len(value.Content)
	*c = make(Cases, 0, l)
	for i := 0; i < l; i += 2 {
		*c = append(*c, value.Content[i].Value, value.Content[i+1].Value)
	}

	return nil
}

func (c *Cases) UnmarshalJSON(data []byte) error {
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
