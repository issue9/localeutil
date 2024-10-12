// SPDX-FileCopyrightText: 2020-2024 caixw
//
// SPDX-License-Identifier: MIT

// Package message 本地化信息的定义
package message

import (
	"slices"
	"strconv"

	"github.com/issue9/sliceutil"
	"golang.org/x/text/feature/plural"
	"golang.org/x/text/language"
	"golang.org/x/text/message/catalog"

	"github.com/issue9/localeutil"
)

type (
	// Language 某一语言的本地化内容
	Language struct {
		XMLName  struct{}     `xml:"language" json:"-" yaml:"-"`
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

	LogFunc = func(localeutil.Stringer)
)

// Join 将 l2.Messages 并入 l.Messages
//
// 执行以下操作：
//
//	-如果 l2 的 [Message.Key] 存在于 l，则覆盖 l 的项；
//	-如果 l2 的 [Message.Key] 不存在于 l，则写入 l；
func (l *Language) Join(l2 *Language) {
	for index, m2 := range l2.Messages {
		elem, found := sliceutil.At(l.Messages, func(m1 Message, _ int) bool { return m1.Key == m2.Key })
		if !found {
			l.Messages = append(l.Messages, m2)
		} else {
			l2.Messages[index] = elem
		}
	}
}

// Merge 将 l.Messages 写入 dest 中的每个元素
//
// 这将会执行以下几个步骤：
//
//	-删除只存在于 dest 元素中而不存在于 l 的内容；
//	-将 l 独有的项写入 dest；
//
// 最终内容是 dest 为准。
// log 所有删除的记录都将通过此输出；
func (l *Language) MergeTo(log LogFunc, dest []*Language) {
	for _, d := range dest {
		l.mergeTo(log, d)
	}
}

func (l *Language) mergeTo(log LogFunc, dest *Language) {
	// 删除只存在于 dest 而不存在于 l 的内容
	dest.Messages = sliceutil.Delete(dest.Messages, func(dm Message, _ int) bool {
		exist := slices.IndexFunc(l.Messages, func(sm Message) bool { return sm.Key == dm.Key }) >= 0
		if !exist {
			log(localeutil.Phrase("the key %s of %s not found, will be deleted", strconv.Quote(dm.Key), dest.ID))
		}
		return !exist
	})

	// 将 l 独有的项写入 dest
	for _, sm := range l.Messages {
		if slices.IndexFunc(dest.Messages, func(dm Message) bool { return dm.Key == sm.Key }) < 0 {
			dest.Messages = append(dest.Messages, sm)
		}
	}
}

// Catalog 将本地化信息附加在 [catalog.Catalog] 上
func (l *Language) Catalog(b *catalog.Builder) (err error) {
	for _, msg := range l.Messages {
		switch {
		case msg.Message.Vars != nil:
			vars := msg.Message.Vars
			msgs := make([]catalog.Message, 0, len(vars))
			for _, v := range vars {
				mm := catalog.Var(v.Name, plural.Selectf(v.Arg, v.Format, ex(v.Cases)...))
				msgs = append(msgs, mm)
			}
			msgs = append(msgs, catalog.String(msg.Message.Msg))
			err = b.Set(l.ID, msg.Key, msgs...)
		case msg.Message.Select != nil:
			s := msg.Message.Select
			err = b.Set(l.ID, msg.Key, plural.Selectf(s.Arg, s.Format, ex(s.Cases)...))
		case msg.Message.Msg != "":
			err = b.SetString(l.ID, msg.Key, msg.Message.Msg)
		}

		if err != nil {
			return err
		}
	}

	return nil
}

func ex(cases []*Case) []any {
	data := make([]any, 0, len(cases)*2)
	for _, c := range cases {
		data = append(data, c.Case, c.Value)
	}
	return data
}
