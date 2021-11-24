// SPDX-License-Identifier: MIT

package localeutil

import (
	"testing"

	"github.com/issue9/assert/v2"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

var (
	_ LocaleStringer = phrase{}
	_ LocaleStringer = &phrase{}

	_ error          = &localeError{}
	_ LocaleStringer = localeError{}
)

func TestLocaleStringer(t *testing.T) {
	a := assert.New(t, false)

	a.NotError(message.SetString(language.SimplifiedChinese, "k1", "cn"))
	a.NotError(message.SetString(language.TraditionalChinese, "k1", "tw"))
	a.NotError(message.SetString(language.SimplifiedChinese, "k2", "cn %[1]s"))
	a.NotError(message.SetString(language.TraditionalChinese, "k2", "tw %[1]s"))
	cnp := message.NewPrinter(language.SimplifiedChinese, message.Catalog(message.DefaultCatalog))
	twp := message.NewPrinter(language.TraditionalChinese, message.Catalog(message.DefaultCatalog))

	p := Phrase("k1")
	a.Equal(p.LocaleString(cnp), "cn")
	a.Equal(p.LocaleString(twp), "tw")
	a.Equal(p.(phrase).String(), "k1")

	p = Phrase("k2", p)
	a.Equal(p.LocaleString(cnp), "cn cn")
	a.Equal(p.LocaleString(twp), "tw tw")
	a.Equal(p.(phrase).String(), "k2")

	p = Phrase("not-exists")
	a.Equal(p.LocaleString(twp), "not-exists")
}

func TestError(t *testing.T) {
	a := assert.New(t, false)

	a.NotError(message.SetString(language.SimplifiedChinese, "k1", "cn"))
	a.NotError(message.SetString(language.TraditionalChinese, "k1", "tw"))
	cnp := message.NewPrinter(language.SimplifiedChinese, message.Catalog(message.DefaultCatalog))
	twp := message.NewPrinter(language.TraditionalChinese, message.Catalog(message.DefaultCatalog))

	err := Error("k1")
	le, ok := err.(LocaleStringer)
	a.True(ok).NotNil(le)
	a.Equal(le.LocaleString(cnp), "cn")
	a.Equal(le.LocaleString(twp), "tw")
	a.Equal(err.Error(), "k1")

	err = Error("not-exists")
	le, ok = err.(LocaleStringer)
	a.True(ok).NotNil(le)
	a.Equal(le.LocaleString(twp), "not-exists")
}
