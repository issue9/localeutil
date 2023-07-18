// SPDX-License-Identifier: MIT

package extract

import (
	"context"
	"log"
	"testing"

	"github.com/issue9/assert/v3"
	"github.com/issue9/localeutil/message"
	"github.com/issue9/sliceutil"
)

func TestExtract(t *testing.T) {
	a := assert.New(t, false)

	msg, err := Extract(context.Background(), "zh-CN", "./testdata", true, log.Default(), "github.com/issue9/localeutil.Phrase")
	a.NotError(err).NotNil(msg).
		NotNil(msg.Languages[0]).
		Equal(msg.Languages[0].ID.String(), "zh-CN")

	m := msg.Languages[0].Messages
	a.NotNil(m).
		Length(sliceutil.Dup(m, func(m1, m2 message.Message) bool { return m1.Key == m2.Key }), 0). // 没有重复值
		Length(m, 6)

	for _, mm := range m {
		t.Log(mm.Key)
	}

	// 添加了 localeutil.Error

	msg, err = Extract(context.Background(), "zh-CN", "./testdata", true, log.Default(), "github.com/issue9/localeutil.Phrase", "github.com/issue9/localeutil.Error")
	a.NotError(err).NotNil(msg).
		NotNil(msg.Languages[0]).
		Equal(msg.Languages[0].ID.String(), "zh-CN")

	m = msg.Languages[0].Messages
	a.NotNil(m).
		Length(sliceutil.Dup(m, func(m1, m2 message.Message) bool { return m1.Key == m2.Key }), 0). // 没有重复值
		Length(m, 8)

	for _, mm := range m {
		t.Log(mm.Key)
	}
}