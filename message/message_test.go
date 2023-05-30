// SPDX-License-Identifier: MIT

package message

import (
	"encoding/json"
	"encoding/xml"
	"os"
	"testing"

	"github.com/issue9/assert/v3"
)

func TestMessages_Load(t *testing.T) {
	a := assert.New(t, false)

	m := &Messages{}

	a.NotError(m.LoadFS(os.DirFS("./testdata"), "cmn-hans.json", json.Unmarshal))
	a.Length(m.Languages, 1).
		Equal(m.Languages[0].ID.String(), "cmn-Hans")

	a.NotError(m.LoadFSGlob(os.DirFS("./testdata"), "*.json", json.Unmarshal))
	a.Length(m.Languages, 1).
		Equal(m.Languages[0].ID.String(), "cmn-Hans")

	a.NotError(m.LoadGlob("./testdata/*.xml", xml.Unmarshal))
	a.Length(m.Languages, 2)
}

func TestMessages_SaveFile(t *testing.T) {
	a := assert.New(t, false)

	m := &Messages{}
	a.NotError(m.LoadFS(os.DirFS("./testdata"), "cmn-hans.json", json.Unmarshal))
	a.NotError(m.LoadFS(os.DirFS("./testdata"), "cmn-hant.xml", xml.Unmarshal))
	a.NotError(m.SaveFiles("./testdata/", ".out", json.Marshal, os.ModePerm))
}
