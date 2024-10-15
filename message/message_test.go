// SPDX-FileCopyrightText: 2020-2024 caixw
//
// SPDX-License-Identifier: MIT

package message

import (
	"testing"

	"github.com/issue9/assert/v4"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"golang.org/x/text/message/catalog"

	"github.com/issue9/localeutil"
)

func TestLanguage_Join(t *testing.T) {
	a := assert.New(t, false)

	src := &File{
		ID:       language.SimplifiedChinese,
		Messages: []Message{{Key: "src"}, {Key: "g", Message: Text{Msg: "src"}}},
	}
	l := &File{
		ID:       language.SimplifiedChinese,
		Messages: []Message{{Key: "l"}, {Key: "g", Message: Text{Msg: "l"}}},
	}
	l.Join(src)
	a.Length(l.Messages, 3)
}

func TestLanguage_MergeTo(t *testing.T) {
	a := assert.New(t, false)
	log := func(s localeutil.Stringer) {}

	dest := &File{
		ID:       language.SimplifiedChinese,
		Messages: []Message{{Key: "dest"}},
	}
	l := &File{
		ID:       language.Afrikaans,
		Messages: []Message{{Key: "l"}},
	}
	l.MergeTo(log, []*File{dest})
	a.Equal(dest.ID, language.SimplifiedChinese).
		Length(dest.Messages, 1).Equal(dest.Messages[0].Key, "l").
		Length(l.Messages, 1).Equal(l.Messages[0].Key, "l")

	dest = &File{
		ID:       language.SimplifiedChinese,
		Messages: []Message{{Key: "dest"}, {Key: "g"}},
	}
	l = &File{
		ID:       language.SimplifiedChinese,
		Messages: []Message{{Key: "l"}, {Key: "g"}},
	}
	l.MergeTo(log, []*File{dest})
	a.Length(dest.Messages, 2).
		Length(l.Messages, 2).Equal(l.Messages[0].Key, "l").Equal(l.Messages[1].Key, "g")
}

func TestLanguage_Catalog(t *testing.T) {
	a := assert.New(t, false)

	b := catalog.NewBuilder()
	l := &File{
		ID: language.SimplifiedChinese,
		Messages: []Message{
			{Key: "k1", Message: Text{Msg: "msg1"}},
			{Key: "k2 %s", Message: Text{Msg: "msg-%s"}},
			{Key: "k3", Message: Text{Select: &Select{
				Arg:    1,
				Format: "%d",
				Cases: []*Case{
					{Case: "=1", Value: "msg-1"},
					{Case: "=2", Value: "msg-2"},
					{Case: "other", Value: "msg-other"},
				},
			}}},
			{Key: "k4", Message: Text{Msg: "${n}-${s}", Vars: []*Var{
				{
					Name:   "n",
					Arg:    2,
					Format: "%d",
					Cases: []*Case{
						{Case: "=1", Value: "1"},
						{Case: "other", Value: "o"},
					},
				},
				{
					Name:   "s",
					Arg:    1,
					Format: "%d",
					Cases: []*Case{
						{Case: "=1", Value: "s1"},
						{Case: "other", Value: "so"},
					},
				},
			}}},
		},
	}
	a.NotError(l.Catalog(b))

	cnp := message.NewPrinter(language.SimplifiedChinese, message.Catalog(b))
	a.Equal(cnp.Sprintf("k1"), "msg1")
	a.Equal(cnp.Sprintf("k2 %s", "1"), "msg-1")
	a.Equal(cnp.Sprintf("k3", 3), "msg-other")
	a.Equal(cnp.Sprintf("k3", 2), "msg-2")
	a.Equal(cnp.Sprintf("k4", 1, 2), "o-s1")
	a.Equal(cnp.Sprintf("k4", 2, 1), "1-so")

	// 未定义 und，cmn-hans 无法找到匹配的数据
	cnp = message.NewPrinter(language.MustParse("cmn-hans"), message.Catalog(b))
	a.Equal(cnp.Sprintf("k1"), "k1")
}
