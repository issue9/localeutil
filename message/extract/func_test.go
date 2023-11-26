// SPDX-License-Identifier: MIT

package extract

import (
	"go/parser"
	"go/token"
	"testing"

	"github.com/issue9/assert/v3"
)

func TestSplit(t *testing.T) {
	a := assert.New(t, false)

	fns := split("github.com/issue9/localeutil.Phrase", "github.com/issue9/localeutil.Error", "github.com/issue9/localeutil.Struct.Printf")
	a.Equal(fns, []importFunc{
		{modName: "github.com/issue9/localeutil", name: "Phrase"},
		{modName: "github.com/issue9/localeutil", name: "Error"},
		{modName: "github.com/issue9/localeutil", name: "Printf", structName: "Struct"},
	})

	a.PanicString(func() {
		split("github.com/issue9")
	}, "github.com/issue9 格式无效")
}

func TestFilterImportFuncs(t *testing.T) {
	t.Run("normal", func(t *testing.T) {
		a := assert.New(t, false)
		f, err := parser.ParseFile(token.NewFileSet(), "./testdata/testdata.go", nil, parser.AllErrors)
		a.NotError(err).NotNil(f)

		fns := []string{"github.com/issue9/localeutil.Phrase", "github.com/issue9/localeutil.Error"}
		mods := filterImportFuncs("", f.Imports, fns)
		a.Equal(mods, []importFunc{
			{modName: "localeutil", name: "Phrase"},
			{modName: "l", name: "Phrase"},
			{modName: "localeutil", name: "Error"},
			{modName: "l", name: "Error"},
		})
	})

	t.Run("signal func", func(t *testing.T) {
		a := assert.New(t, false)
		f, err := parser.ParseFile(token.NewFileSet(), "./testdata/testdata.go", nil, parser.AllErrors)
		a.NotError(err).NotNil(f)

		fns := []string{"github.com/issue9/localeutil.Phrase"}
		mods := filterImportFuncs("", f.Imports, fns)
		a.Equal(mods, []importFunc{
			{modName: "localeutil", name: "Phrase"},
			{modName: "l", name: "Phrase"},
		})
	})

	t.Run("struct", func(t *testing.T) {
		a := assert.New(t, false)
		f, err := parser.ParseFile(token.NewFileSet(), "./testdata/struct.go", nil, parser.AllErrors)
		a.NotError(err).NotNil(f)

		fns := []string{"golang.org/x/text/message.Printer.Printf"}
		mods := filterImportFuncs("", f.Imports, fns)
		a.Equal(mods, []importFunc{
			{modName: "message", name: "Printf", structName: "Printer"},
			{modName: "xm", name: "Printf", structName: "Printer"},
		})
	})

	t.Run("file mod path", func(t *testing.T) {
		a := assert.New(t, false)
		f, err := parser.ParseFile(token.NewFileSet(), "./testdata/struct.go", nil, parser.AllErrors)
		a.NotError(err).NotNil(f)

		fns := []string{"golang.org/x/text/message.Printer.Printf"}
		mods := filterImportFuncs("golang.org/x/text/message", f.Imports, fns)
		a.Equal(mods, []importFunc{
			{name: "Printf", structName: "Printer"},
		})
	})
}
