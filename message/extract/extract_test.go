// SPDX-License-Identifier: MIT

package extract

import (
	"context"
	"log"
	"testing"

	"github.com/issue9/assert/v3"
	"github.com/issue9/sliceutil"
	"golang.org/x/text/language"
	xm "golang.org/x/text/message"
	"golang.org/x/text/message/catalog"

	"github.com/issue9/localeutil"
	"github.com/issue9/localeutil/message"
)

func TestExtract(t *testing.T) {
	a := assert.New(t, false)
	b := catalog.NewBuilder()
	p := xm.NewPrinter(language.SimplifiedChinese, xm.Catalog(b))
	log := func(v localeutil.Stringer) { log.Print(v.LocaleString(p)) }

	o := &Options{
		Language:  language.MustParse("zh-CN"),
		Root:      "./testdata",
		Recursive: true,
		Log:       log,
		Funcs:     []string{"github.com/issue9/localeutil.Phrase"},
	}
	l, err := Extract(context.Background(), o)
	a.NotError(err).NotNil(l).
		Equal(l.ID.String(), "zh-CN")

	m := l.Messages
	a.NotNil(m).
		Length(sliceutil.Dup(m, func(m1, m2 message.Message) bool { return m1.Key == m2.Key }), 0). // 没有重复值
		Length(m, 5).
		Equal(m[0].Key, "const-value")

	for _, mm := range m {
		t.Log(mm.Key)
	}

	// 添加了 localeutil.Error 和 localeutil.StringPhrase

	o = &Options{
		Language:  language.MustParse("zh-CN"),
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
		NotNil(l).
		Equal(l.ID.String(), "zh-CN")

	m = l.Messages
	a.NotNil(m).
		Length(sliceutil.Dup(m, func(m1, m2 message.Message) bool { return m1.Key == m2.Key }), 0). // 没有重复值
		Length(m, 9)

	for _, mm := range m {
		t.Log(mm.Key)
	}

	// 添加了 text/message.Printer.Printf

	o = &Options{
		Language:  language.MustParse("zh-CN"),
		Root:      "./testdata",
		Recursive: true,
		Log:       log,
		Funcs: []string{
			"github.com/issue9/localeutil.Phrase",
			"github.com/issue9/localeutil.Error",
			"github.com/issue9/localeutil.StringPhrase",
			"golang.org/x/text/message.Printer.Printf",
		},
	}
	l, err = Extract(context.Background(), o)
	a.NotError(err).NotNil(l).
		NotNil(l).
		Equal(l.ID.String(), "zh-CN")

	m = l.Messages
	a.NotNil(m).
		Length(m, 12)

	for _, mm := range m {
		t.Log(mm.Key)
	}

	// 测试本地的函数和对象

	o = &Options{
		Language:  language.MustParse("zh-CN"),
		Root:      "./testdata",
		Recursive: true,
		Log:       log,
		Funcs: []string{
			"github.com/issue9/localeutil.Phrase",
			"github.com/issue9/localeutil.Error",
			"github.com/issue9/localeutil.StringPhrase",
			"golang.org/x/text/message.Printer.Printf",
			"github.com/issue9/localeutil/testdata.Printer.Print",
			"github.com/issue9/localeutil/testdata.Print",
		},
	}
	l, err = Extract(context.Background(), o)
	a.NotError(err).NotNil(l).
		NotNil(l).
		Equal(l.ID.String(), "zh-CN")

	m = l.Messages
	a.NotNil(m).
		Length(m, 14)

	for _, mm := range m {
		t.Log(mm.Key)
	}
}
