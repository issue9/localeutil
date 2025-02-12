// SPDX-FileCopyrightText: 2020-2025 caixw
//
// SPDX-License-Identifier: MIT

package serialize

import (
	"encoding/json"
	"encoding/xml"
	"os"
	"testing"

	"github.com/issue9/assert/v4"
	"golang.org/x/text/language"
)

func TestLoad(t *testing.T) {
	a := assert.New(t, false)

	l, err := LoadFS(os.DirFS("./testdata"), "cmn-hans.json", json.Unmarshal)
	a.NotError(err).NotNil(l)
	a.Equal(l.Languages, []language.Tag{language.MustParse("cmn-Hans")})

	ls, err := LoadFSGlob(func(string) UnmarshalFunc { return json.Unmarshal }, "*.json", os.DirFS("./testdata"))
	a.NotError(err).Length(ls, 1)
	a.Length(ls, 1).
		Equal(ls[0].Languages, []language.Tag{language.MustParse("cmn-Hans")})

	ls, err = LoadGlob(func(string) UnmarshalFunc { return xml.Unmarshal }, "./testdata/*.xml")
	a.NotError(err).Length(ls, 1)
}
