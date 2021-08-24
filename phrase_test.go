// SPDX-License-Identifier: MIT

package localeutil

import (
	"testing"

	"github.com/issue9/assert"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

func TestPhrase_LocaleString(t *testing.T) {
	a := assert.New(t)

	message.SetString(language.SimplifiedChinese, "k1", "cn")
	message.SetString(language.TraditionalChinese, "k1", "tw")
	cnp := message.NewPrinter(language.SimplifiedChinese, message.Catalog(message.DefaultCatalog))
	twp := message.NewPrinter(language.TraditionalChinese, message.Catalog(message.DefaultCatalog))

	p := Phrase{Key: "k1"}
	a.Equal(p.LocaleString(cnp), "cn")
	a.Equal(p.LocaleString(twp), "tw")
}
