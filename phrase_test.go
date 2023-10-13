// SPDX-License-Identifier: MIT

package localeutil

import (
	"errors"
	"fmt"
	"testing"

	"github.com/issue9/assert/v3"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

var (
	_ Stringer = phrase{}
	_ Stringer = &phrase{}
	_ Stringer = StringPhrase("123")

	_ error    = &localeError{}
	_ Stringer = &localeError{}
)

func TestStringer(t *testing.T) {
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

	p = Phrase("k2", p)
	a.Equal(p.LocaleString(cnp), "cn cn")
	a.Equal(p.LocaleString(twp), "tw tw")

	p = Phrase("not-exists")
	a.Equal(p.LocaleString(twp), "not-exists")

	p = StringPhrase("not-exists")
	a.Equal(p.LocaleString(twp), "not-exists")

	p = StringPhrase("k1")
	a.Equal(p.LocaleString(cnp), "cn")
	a.Equal(p.LocaleString(twp), "tw")
}

func TestError(t *testing.T) {
	a := assert.New(t, false)

	a.NotError(message.SetString(language.SimplifiedChinese, "k1", "cn"))
	a.NotError(message.SetString(language.TraditionalChinese, "k1", "tw"))
	cnp := message.NewPrinter(language.SimplifiedChinese, message.Catalog(message.DefaultCatalog))
	twp := message.NewPrinter(language.TraditionalChinese, message.Catalog(message.DefaultCatalog))

	err := Error("k1")
	le, ok := err.(Stringer)
	a.True(ok).NotNil(le)
	a.Equal(le.LocaleString(cnp), "cn")
	a.Equal(le.LocaleString(twp), "tw")
	a.Equal(err.Error(), "k1")

	err = Error("not-exists")
	le, ok = err.(Stringer)
	a.True(ok).NotNil(le)
	a.Equal(le.LocaleString(twp), "not-exists")

	err1 := Error("k1")
	a.True(errors.Is(err1, err1))
	a.True(errors.Is(fmt.Errorf("err2 %w", err1), err1))
	a.False(errors.Is(Error("k1"), Error("k1")))
	a.False(errors.Is(Error("k1"), errors.New("k1")))

	// ErrorAsLocaleString

	a.Equal(ErrorAsLocaleString(err1, cnp), "cn")
	a.Equal(ErrorAsLocaleString(err1, twp), "tw")
	a.Equal(ErrorAsLocaleString(errors.New("k1"), twp), "k1")
}

func BenchmarkPhrase_LocaleString(b *testing.B) {
	a := assert.New(b, false)

	a.NotError(message.SetString(language.SimplifiedChinese, "k1", "cn"))
	a.NotError(message.SetString(language.SimplifiedChinese, "k2", "cn %[1]s"))
	a.NotError(message.SetString(language.SimplifiedChinese, "k3", "cn %[1]s"))
	a.NotError(message.SetString(language.SimplifiedChinese, "k4", "cn %[1]s %[2]s %[3]s"))
	cnp := message.NewPrinter(language.SimplifiedChinese, message.Catalog(message.DefaultCatalog))

	p1 := Phrase("k1")
	b.Run("0 Stringer", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			a.Equal(p1.LocaleString(cnp), "cn")
		}
	})

	p2 := Phrase("k2", p1)
	b.Run("1 Stringer", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			a.Equal(p2.LocaleString(cnp), "cn cn")
		}
	})

	p3 := Phrase("k3", p2)
	b.Run("2 Stringer", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			a.Equal(p3.LocaleString(cnp), "cn cn cn")
		}
	})

	p4 := Phrase("k4", p1, p1, p1)
	b.Run("3 Stringer", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			a.Equal(p4.LocaleString(cnp), "cn cn cn cn")
		}
	})
}
