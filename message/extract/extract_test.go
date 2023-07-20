// SPDX-License-Identifier: MIT

package extract

import (
	"context"
	"log"
	"testing"

	"github.com/issue9/assert/v3"
	"github.com/issue9/sliceutil"

	"github.com/issue9/localeutil/message"
)

func TestExtract(t *testing.T) {
	a := assert.New(t, false)

	o := &Options{
		Language:  "zh-CN",
		Root:      "./testdata",
		Recursive: true,
		Log:       log.Default(),
		Funcs:     []string{"github.com/issue9/localeutil.Phrase"},
	}
	msg, err := Extract(context.Background(), o)
	a.NotError(err).NotNil(msg).
		NotNil(msg.Languages[0]).
		Equal(msg.Languages[0].ID.String(), "zh-CN")

	m := msg.Languages[0].Messages
	a.NotNil(m).
		Length(sliceutil.Dup(m, func(m1, m2 message.Message) bool { return m1.Key == m2.Key }), 0). // 没有重复值
		Length(m, 5).
		Equal(m[0].Key, "const-value")

	for _, mm := range m {
		t.Log(mm.Key)
	}

	// 添加了 localeutil.Error 和 localeutil.StringPhrase

	o = &Options{
		Language:  "zh-CN",
		Root:      "./testdata",
		Recursive: true,
		Log:       log.Default(),
		Funcs: []string{
			"github.com/issue9/localeutil.Phrase",
			"github.com/issue9/localeutil.Error",
			"github.com/issue9/localeutil.StringPhrase",
		},
	}
	msg, err = Extract(context.Background(), o)
	a.NotError(err).NotNil(msg).
		NotNil(msg.Languages[0]).
		Equal(msg.Languages[0].ID.String(), "zh-CN")

	m = msg.Languages[0].Messages
	a.NotNil(m).
		Length(sliceutil.Dup(m, func(m1, m2 message.Message) bool { return m1.Key == m2.Key }), 0). // 没有重复值
		Length(m, 9)

	for _, mm := range m {
		t.Log(mm.Key)
	}

	// 添加了 text/message.Printer.Printf

	o = &Options{
		Language:  "zh-CN",
		Root:      "./testdata",
		Recursive: true,
		Log:       log.Default(),
		Funcs: []string{
			"github.com/issue9/localeutil.Phrase",
			"github.com/issue9/localeutil.Error",
			"github.com/issue9/localeutil.StringPhrase",
			"golang.org/x/text/message.Printer.Printf",
		},
	}
	msg, err = Extract(context.Background(), o)
	a.NotError(err).NotNil(msg).
		NotNil(msg.Languages[0]).
		Equal(msg.Languages[0].ID.String(), "zh-CN")

	m = msg.Languages[0].Messages
	a.NotNil(m).
		Length(m, 12)

	for _, mm := range m {
		t.Log(mm.Key)
	}

	// 测试本地的函数和对象

	o = &Options{
		Language:  "zh-CN",
		Root:      "./testdata",
		Recursive: true,
		Log:       log.Default(),
		Funcs: []string{
			"github.com/issue9/localeutil.Phrase",
			"github.com/issue9/localeutil.Error",
			"github.com/issue9/localeutil.StringPhrase",
			"golang.org/x/text/message.Printer.Printf",
			"github.com/issue9/localeutil/testdata.Printer.Print",
			"github.com/issue9/localeutil/testdata.Print",
		},
	}
	msg, err = Extract(context.Background(), o)
	a.NotError(err).NotNil(msg).
		NotNil(msg.Languages[0]).
		Equal(msg.Languages[0].ID.String(), "zh-CN")

	m = msg.Languages[0].Messages
	a.NotNil(m).
		Length(m, 14)

	for _, mm := range m {
		t.Log(mm.Key)
	}
}
