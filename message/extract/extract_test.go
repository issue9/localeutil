// SPDX-FileCopyrightText: 2020-2024 caixw
//
// SPDX-License-Identifier: MIT

package extract

import (
	"context"
	"log"
	"testing"

	"github.com/issue9/assert/v4"
	"github.com/issue9/sliceutil"
	"golang.org/x/text/language"

	"github.com/issue9/localeutil"
	"github.com/issue9/localeutil/message"
)

func TestExtract(t *testing.T) {
	a := assert.New(t, false)
	log := func(v localeutil.Stringer) { log.Print(v.LocaleString(nil)) }

	// 全部是本地对象

	t.Run("testdata/locale", func(t *testing.T) {
		o := &Options{
			Language:  language.MustParse("zh-CN"),
			Root:      "./testdata/locale",
			Recursive: true,
			WarnLog:   log,
			InfoLog:   log,
			Funcs: []string{
				"github.com/issue9/localeutil/testdata/locale.String",
				"github.com/issue9/localeutil/testdata/locale.Print",
				"github.com/issue9/localeutil/testdata/locale.GPrinter.GPrint",
				"github.com/issue9/localeutil/testdata/locale.Printer.Print",
				"github.com/issue9/localeutil/testdata/locale.Interface.Printf",
			},
		}
		l, err := Extract(context.Background(), o)
		a.NotError(err).NotNil(l).
			Equal(l.Languages, []language.Tag{language.MustParse("zh-CN")})

		m := l.Messages
		a.NotNil(m).
			Length(sliceutil.Dup(m, func(m1, m2 message.Message) bool { return m1.Key == m2.Key }), 0). // 没有重复值
			Length(m, 14)

		for _, mm := range m {
			t.Log(mm.Key)
		}
		t.Log("\n\n")
	})

	// 单个 localeutil.Phrase
	t.Run("localeutilx1", func(t *testing.T) {
		o := &Options{
			Root:    "./testdata",
			WarnLog: log,
			Funcs:   []string{"github.com/issue9/localeutil.Phrase"},
		}
		l, err := Extract(context.Background(), o)
		a.NotError(err).NotNil(l)

		m := l.Messages
		a.Length(m, 5).
			Length(sliceutil.Dup(m, func(m1, m2 message.Message) bool { return m1.Key == m2.Key }), 0) // 没有重复值

		for _, mm := range m {
			t.Log(mm.Key)
		}
		t.Log("\n\n")
	})

	// 添加了 localeutil.Error 和 localeutil.StringPhrase，
	t.Run("localeutilx3", func(t *testing.T) {
		o := &Options{
			Root:      "./testdata",
			Recursive: true,
			WarnLog:   log,
			Funcs: []string{
				"github.com/issue9/localeutil.Phrase",
				"github.com/issue9/localeutil.Error",
				"github.com/issue9/localeutil.StringPhrase",
				//"github.com/issue9/localeutil/testdata.String", // 别名，指向 localeutil.StringPhrase
			},
		}
		l, err := Extract(context.Background(), o)
		a.NotError(err).NotNil(l).
			NotNil(l)

		m := l.Messages
		a.Length(m, 10).
			Length(sliceutil.Dup(m, func(m1, m2 message.Message) bool { return m1.Key == m2.Key }), 0) // 没有重复值

		for _, mm := range m {
			t.Log(mm.Key)
		}
		t.Log("\n\n")
	})

	// 添加了 locale.Printer.Print，以及 struct tag
	t.Run("localeutilx3,localex2", func(t *testing.T) {
		o := &Options{
			Root:      "./testdata",
			Recursive: true,
			WarnLog:   log,
			InfoLog:   log,
			Funcs: []string{
				"github.com/issue9/localeutil.Phrase",
				"github.com/issue9/localeutil.Error",
				"github.com/issue9/localeutil.StringPhrase",
				"github.com/issue9/localeutil/testdata/locale.String",        // 别名，ref.String 指向此值
				"github.com/issue9/localeutil/testdata/locale.Printer.Print", // ref.Printer 指向此值
			},
			Tag: "comment",
		}
		l, err := Extract(context.Background(), o)
		a.NotError(err).NotNil(l).
			NotNil(l)

		m := l.Messages
		a.Length(m, 24)

		for _, mm := range m {
			t.Log(mm.Key)
		}
		t.Log("\n\n")
	})
}

func TestParseTypeName(t *testing.T) {
	a := assert.New(t, false)

	p, s := parseTypeName("t")
	a.Empty(p).Equal(s, "t")

	p, s = parseTypeName("*github.com/issue9/abc.t")
	a.Equal(p, "github.com/issue9/abc").Equal(s, "t")

	p, s = parseTypeName("github.com/issue9/abc.t[int]")
	a.Equal(p, "github.com/issue9/abc").Equal(s, "t")

	p, s = parseTypeName("*github.com/issue9/abc.t[github.com/issue9/abc.Type]")
	a.Equal(p, "github.com/issue9/abc").Equal(s, "t")

	p, s = parseTypeName("github.com/issue9/abc.t[github.com/issue9/abc.Type]")
	a.Equal(p, "github.com/issue9/abc").Equal(s, "t")
}
