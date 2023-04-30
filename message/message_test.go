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
	"gopkg.in/yaml.v3"
)

var (
	_ yaml.Unmarshaler = &Cases{}
	_ xml.Unmarshaler  = &Cases{}
	_ json.Unmarshaler = &Cases{}
)

func TestLoadFromFS_yaml(t *testing.T) {
	a := assert.New(t, false)
	b := catalog.NewBuilder()

	a.NotError(LoadFromFS(b, os.DirFS("./testdata"), "cmn-hans.yaml", yaml.Unmarshal))
	p := message.NewPrinter(language.MustParse("cmn-hans"), message.Catalog(b))

	a.Equal(p.Sprintf("k1"), "msg1")

	a.Equal(p.Sprintf("k2", 1), "msg-1")
	a.Equal(p.Sprintf("k2", 3), "msg-3")
	a.Equal(p.Sprintf("k2", 5), "msg-other")

	a.Equal(p.Sprintf("k3", 1, 1), "1-一")
	a.Equal(p.Sprintf("k3", 1, 2), "2-一")
	a.Equal(p.Sprintf("k3", 2, 2), "2-二")

	// yaml 中也定义了 und
	p = message.NewPrinter(language.Albanian, message.Catalog(b))
	a.Equal(p.Sprintf("k1"), "msg1")
}

func TestLoadFromFSGlob_yaml(t *testing.T) {
	a := assert.New(t, false)
	b := catalog.NewBuilder()

	a.NotError(LoadFromFSGlob(b, os.DirFS("./testdata"), "*.yaml", yaml.Unmarshal))
	p := message.NewPrinter(language.MustParse("cmn-hans"), message.Catalog(b))

	a.Equal(p.Sprintf("k1"), "msg1")

	a.Equal(p.Sprintf("k2", 1), "msg-1")
	a.Equal(p.Sprintf("k2", 3), "msg-3")
	a.Equal(p.Sprintf("k2", 5), "msg-other")

	a.Equal(p.Sprintf("k3", 1, 1), "1-一")
	a.Equal(p.Sprintf("k3", 1, 2), "2-一")
	a.Equal(p.Sprintf("k3", 2, 2), "2-二")

	// yaml 中也定义了 und
	p = message.NewPrinter(language.Albanian, message.Catalog(b))
	a.Equal(p.Sprintf("k1"), "msg1")
}

func TestLoadFromFS_xml(t *testing.T) {
	a := assert.New(t, false)
	b := catalog.NewBuilder()

	a.NotError(LoadFromFS(b, os.DirFS("./testdata"), "cmn-hant.xml", xml.Unmarshal))
	p := message.NewPrinter(language.MustParse("cmn-hant"), message.Catalog(b))

	a.Equal(p.Sprintf("k1"), "msg1")

	a.Equal(p.Sprintf("k2", 1), "msg-1")
	a.Equal(p.Sprintf("k2", 3), "msg-3")
	a.Equal(p.Sprintf("k2", 5), "msg-other")

	a.Equal(p.Sprintf("k3", 1, 1), "1-一")
	a.Equal(p.Sprintf("k3", 1, 2), "2-一")
	a.Equal(p.Sprintf("k3", 2, 2), "2-二")

	// 未定义 und，cmn-hans 无法找到匹配的数据
	p = message.NewPrinter(language.MustParse("cmn-hans"), message.Catalog(b))
	a.Equal(p.Sprintf("k1"), "k1")
}

func TestLoadFromFS_json(t *testing.T) {
	a := assert.New(t, false)
	b := catalog.NewBuilder()

	a.NotError(LoadFromFS(b, os.DirFS("./testdata"), "cmn-hans.json", json.Unmarshal))
	p := message.NewPrinter(language.MustParse("cmn-hans"), message.Catalog(b))

	a.Equal(p.Sprintf("k1"), "msg1")

	a.Equal(p.Sprintf("k2", 1), "msg-1")
	a.Equal(p.Sprintf("k2", 3), "msg-3")
	a.Equal(p.Sprintf("k2", 5), "msg-other")

	a.Equal(p.Sprintf("k3", 1, 1), "1-一")
	a.Equal(p.Sprintf("k3", 1, 2), "2-一")
	a.Equal(p.Sprintf("k3", 2, 2), "2-二")
}
