// SPDX-FileCopyrightText: 2025 caixw
//
// SPDX-License-Identifier: MIT

package serialize

import (
	"encoding/json"
	"encoding/xml"
	"io/fs"
	"testing"

	"github.com/issue9/assert/v4"
	"golang.org/x/text/language"

	"github.com/issue9/localeutil/message"
)

func TestSaveFile(t *testing.T) {
	a := assert.New(t, false)

	f := &message.File{
		Languages: []language.Tag{language.SimplifiedChinese},
		Messages: []message.Message{
			{
				Key: "k1",
				Message: message.Text{
					Msg: "m1",
				},
			},
			{
				Key: "k2",
				Message: message.Text{
					Msg: "m2",
				},
			},
		},
	}

	a.NotError(SaveFile(f, "./testdata/json.out", json.Marshal, fs.ModePerm))
	a.NotError(SaveFile(f, "./testdata/xml.out", xml.Marshal, fs.ModePerm))
}
