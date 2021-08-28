// SPDX-License-Identifier: MIT

package localeutil

import (
	"encoding/json"
	"testing"

	"github.com/issue9/assert"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"golang.org/x/text/message/catalog"
)

func TestLoadMessagesFromFile(t *testing.T) {
	a := assert.New(t)

	b := catalog.NewBuilder()
	a.NotError(LoadMessageFromFile(b, "./internal/message/testdata/cmn-hans.json", json.Unmarshal))
	p := message.NewPrinter(language.MustParse("cmn-hans"), message.Catalog(b))
	phrase := Phrase("k3", 3, 1)
	a.Equal(phrase.LocaleString(p), "1-二")
}
