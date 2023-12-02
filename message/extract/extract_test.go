// SPDX-License-Identifier: MIT

package extract

import (
	"context"
	"log"
	"testing"

	"github.com/issue9/assert/v3"
	"github.com/issue9/sliceutil"
	"golang.org/x/text/language"

	"github.com/issue9/localeutil"
	"github.com/issue9/localeutil/message"
)

func TestExtract_LocaleType(t *testing.T) {
	a := assert.New(t, false)
	log := func(v localeutil.Stringer) { log.Print(v.LocaleString(nil)) }

	o := &Options{
		Language:  language.MustParse("zh-CN"),
		Root:      "./testdata/locale",
		Recursive: true,
		Log:       log,
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
		Equal(l.ID.String(), "zh-CN")

	m := l.Messages
	a.NotNil(m).
		Length(sliceutil.Dup(m, func(m1, m2 message.Message) bool { return m1.Key == m2.Key }), 0). // 没有重复值
		Length(m, 13)

	for _, mm := range m {
		t.Log(mm.Key)
	}
}

func TestExtract(t *testing.T) {
	a := assert.New(t, false)
	log := func(v localeutil.Stringer) { log.Print(v.LocaleString(nil)) }

	o := &Options{
		Root:  "./testdata",
		Log:   log,
		Funcs: []string{"github.com/issue9/localeutil.Phrase"},
	}
	l, err := Extract(context.Background(), o)
	a.NotError(err).NotNil(l)

	m := l.Messages
	a.NotNil(m).
		Length(sliceutil.Dup(m, func(m1, m2 message.Message) bool { return m1.Key == m2.Key }), 0). // 没有重复值
		Length(m, 5)

	for _, mm := range m {
		t.Log(mm.Key)
	}

	// 添加了 localeutil.Error 和 localeutil.StringPhrase

	o = &Options{
		Root:      "./testdata",
		Recursive: true,
		Log:       log,
		Funcs: []string{
			"github.com/issue9/localeutil.Phrase",
			"github.com/issue9/localeutil.Error",
			"github.com/issue9/localeutil.StringPhrase",
		},
	}
	l, err = Extract(context.Background(), o)
	a.NotError(err).NotNil(l).
		NotNil(l)

	m = l.Messages
	a.NotNil(m).
		Length(sliceutil.Dup(m, func(m1, m2 message.Message) bool { return m1.Key == m2.Key }), 0). // 没有重复值
		Length(m, 8)

	for _, mm := range m {
		t.Log(mm.Key)
	}

	// 添加了 locale.Printer.Print

	o = &Options{
		Root:      "./testdata",
		Recursive: true,
		Log:       log,
		Funcs: []string{
			"github.com/issue9/localeutil.Phrase",
			"github.com/issue9/localeutil.Error",
			"github.com/issue9/localeutil.StringPhrase",
			"github.com/issue9/localeutil/testdata/locale.String",
		},
	}
	l, err = Extract(context.Background(), o)
	a.NotError(err).NotNil(l).
		NotNil(l)

	m = l.Messages
	a.NotNil(m).
		Length(m, 10)

	for _, mm := range m {
		t.Log(mm.Key)
	}
}

func TestGetTypeName(t *testing.T) {
	a := assert.New(t, false)

	p, s := getTypeName("t")
	a.Empty(p).Equal(s, "t")

	p, s = getTypeName("*github.com/issue9/abc.t")
	a.Equal(p, "github.com/issue9/abc").Equal(s, "t")

	p, s = getTypeName("github.com/issue9/abc.t[int]")
	a.Equal(p, "github.com/issue9/abc").Equal(s, "t")

	p, s = getTypeName("*github.com/issue9/abc.t[github.com/issue9/abc.Type]")
	a.Equal(p, "github.com/issue9/abc").Equal(s, "t")

	p, s = getTypeName("github.com/issue9/abc.t[github.com/issue9/abc.Type]")
	a.Equal(p, "github.com/issue9/abc").Equal(s, "t")
}
