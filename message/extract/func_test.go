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

	fns := split("github.com/issue9/localeutil.Phrase", "github.com/issue9/localeutil.Error")
	a.Equal(fns, []localeFunc{
		{path: "github.com/issue9/localeutil", name: "Phrase"},
		{path: "github.com/issue9/localeutil", name: "Error"},
	})
}

func TestFilterImportFuncs(t *testing.T) {
	t.Run("normal", func(t *testing.T) {
		a := assert.New(t, false)
		f, err := parser.ParseFile(token.NewFileSet(), "./testdata/testdata.go", nil, parser.AllErrors)
		a.NotError(err).NotNil(f)

		fns := split("github.com/issue9/localeutil.Phrase", "github.com/issue9/localeutil.Error")
		mods := filterImportFuncs(f.Imports, fns)
		a.Equal(mods, []importFunc{
			{modName: "localeutil", name: "Phrase"},
			{modName: "localeutil", name: "Error"},
		})
	})

	t.Run("signal func", func(t *testing.T) {
		a := assert.New(t, false)
		f, err := parser.ParseFile(token.NewFileSet(), "./testdata/testdata.go", nil, parser.AllErrors)
		a.NotError(err).NotNil(f)

		fns := split("github.com/issue9/localeutil.Phrase")
		mods := filterImportFuncs(f.Imports, fns)
		a.Equal(mods, []importFunc{
			{modName: "localeutil", name: "Phrase"},
		})
	})

	t.Run("alias mod", func(t *testing.T) {
		a := assert.New(t, false)
		f, err := parser.ParseFile(token.NewFileSet(), "./testdata/alias.go", nil, parser.AllErrors)
		a.NotError(err).NotNil(f)

		fns := split("github.com/issue9/localeutil.Phrase", "github.com/issue9/localeutil.Error")
		mods := filterImportFuncs(f.Imports, fns)
		a.Equal(mods, []importFunc{
			{modName: "localeutil", name: "Phrase"},
			{modName: "localeutil", name: "Error"},
			{modName: "l", name: "Phrase"},
			{modName: "l", name: "Error"},
		})
	})
}
