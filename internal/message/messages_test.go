// SPDX-License-Identifier: MIT

package message

import (
	"encoding/json"
	"encoding/xml"
	"os"
	"testing"

	"github.com/issue9/assert"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"golang.org/x/text/message/catalog"
	"gopkg.in/yaml.v2"
)

func TestLoadFromFS_yaml(t *testing.T) {
	a := assert.New(t)
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
}

func TestLoadFromFS_xml(t *testing.T) {
	a := assert.New(t)
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
}

func TestLoadFromFS_json(t *testing.T) {
	a := assert.New(t)
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
