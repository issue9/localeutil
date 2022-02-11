// SPDX-License-Identifier: MIT

package localeutil

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/issue9/assert/v2"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"golang.org/x/text/message/catalog"
)

func TestLoadMessagesFromFS(t *testing.T) {
	a := assert.New(t, false)

	b := catalog.NewBuilder()
	fsys := os.DirFS("./internal/message/testdata")
	a.NotError(LoadMessageFromFS(b, fsys, "cmn-hans.json", json.Unmarshal))
	p := message.NewPrinter(language.MustParse("cmn-hans"), message.Catalog(b))
	phrase := Phrase("k3", 3, 1)
	a.Equal(phrase.LocaleString(p), "1-二")
}

func TestLoadMessagesFromFSGlob(t *testing.T) {
	a := assert.New(t, false)

	b := catalog.NewBuilder()
	fsys := os.DirFS("./internal/message/testdata")
	a.NotError(LoadMessageFromFSGlob(b, fsys, "*.json", json.Unmarshal))
	p := message.NewPrinter(language.MustParse("cmn-hans"), message.Catalog(b))
	phrase := Phrase("k3", 3, 1)
	a.Equal(phrase.LocaleString(p), "1-二")
}
