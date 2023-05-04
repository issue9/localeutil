// SPDX-License-Identifier: MIT

package message

import (
	"encoding/json"
	"encoding/xml"
	"os"
	"testing"

	"github.com/issue9/assert/v3"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"golang.org/x/text/message/catalog"
)

func TestMessages_Catalog(t *testing.T) {
	t.Run("xml", func(t *testing.T) {
		a := assert.New(t, false)
		b := catalog.NewBuilder()
		m := &Messages{}

		a.NotError(m.LoadFS(os.DirFS("./testdata"), "cmn-hant.xml", xml.Unmarshal))
		m.Catalog(b)
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
		m := &Messages{}

		a.NotError(m.LoadGlob("./testdata/*.json", json.Unmarshal))
		m.Catalog(b)
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
