// SPDX-FileCopyrightText: 2020-2024 caixw
//
// SPDX-License-Identifier: MIT

package serialize

import (
	"encoding/json"
	"encoding/xml"
	"os"
	"testing"

	"github.com/issue9/assert/v4"

	"github.com/issue9/localeutil/message"
)

func TestLoad(t *testing.T) {
	a := assert.New(t, false)

	l, err := LoadFS(os.DirFS("./testdata"), "cmn-hans.json", json.Unmarshal)
	a.NotError(err).NotNil(l)
	a.Equal(l.ID.String(), "cmn-Hans")

	ls, err := LoadFSGlob(func(string) UnmarshalFunc { return json.Unmarshal }, "*.json", os.DirFS("./testdata"))
	a.NotError(err).Length(ls, 1)
	a.Length(ls, 1).
		Equal(ls[0].ID.String(), "cmn-Hans")

	ls, err = LoadGlob(func(string) UnmarshalFunc { return xml.Unmarshal }, "./testdata/*.xml")
	a.NotError(err).Length(ls, 1)
}

func TestSaveFile(t *testing.T) {
	a := assert.New(t, false)

	l1, err := LoadFS(os.DirFS("./testdata"), "cmn-hans.json", json.Unmarshal)
	a.NotError(err).NotNil(l1)
	l2, err := LoadFS(os.DirFS("./testdata"), "cmn-hant.xml", xml.Unmarshal)
	a.NotError(err).NotNil(l2)

	a.NotError(SaveFiles([]*message.Language{l1, l2}, "./testdata/", ".out", json.Marshal, os.ModePerm))
}
