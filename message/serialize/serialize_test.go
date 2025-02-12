// SPDX-FileCopyrightText: 2020-2024 caixw
//
// SPDX-License-Identifier: MIT

package serialize

import (
	"encoding/json"
	"encoding/xml"
	"os"
	"testing"

	"github.com/issue9/assert/v4"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"golang.org/x/text/message/catalog"
)

func TestCatalog(t *testing.T) {
	t.Run("xml", func(t *testing.T) {
		a := assert.New(t, false)
		b := catalog.NewBuilder()

		l, err := LoadFS(os.DirFS("./testdata"), "cmn-hant.xml", xml.Unmarshal)
		a.NotError(err).NotNil(l)
		l.Catalog(b)
		hant := message.NewPrinter(language.MustParse("cmn-hant"), message.Catalog(b))

		a.Equal(hant.Sprintf("k1"), "msg1")

		a.Equal(hant.Sprintf("k2", 1), "msg-1")
		a.Equal(hant.Sprintf("k2", 3), "msg-3")
		a.Equal(hant.Sprintf("k2", 5), "msg-other")

		a.Equal(hant.Sprintf("k3", 1, 1), "1-一")
		a.Equal(hant.Sprintf("k3", 1, 2), "2-一")
		a.Equal(hant.Sprintf("k3", 2, 2), "2-二")

		// 未定义 und，cmn-hans 无法找到匹配的数据
		hant = message.NewPrinter(language.MustParse("cmn-hans"), message.Catalog(b))
		a.Equal(hant.Sprintf("k1"), "k1")
	})

	t.Run("json", func(t *testing.T) {
		a := assert.New(t, false)
		b := catalog.NewBuilder()

		ls, err := LoadGlob(func(string) UnmarshalFunc { return json.Unmarshal }, "./testdata/*.json")
		a.NotError(err).Length(ls, 1)
		for _, l := range ls {
			l.Catalog(b)
		}
		p := message.NewPrinter(language.MustParse("cmn-hans"), message.Catalog(b))

		a.Equal(p.Sprintf("k1"), "msg1")

		a.Equal(p.Sprintf("k2", 1), "msg-1")
		a.Equal(p.Sprintf("k2", 3), "msg-3")
		a.Equal(p.Sprintf("k2", 5), "msg-other")

		a.Equal(p.Sprintf("k3", 1, 1), "1-一")
		a.Equal(p.Sprintf("k3", 1, 2), "2-一")
		a.Equal(p.Sprintf("k3", 2, 2), "2-二")
	})
}
