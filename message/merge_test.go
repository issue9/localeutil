// SPDX-License-Identifier: MIT

package message

import (
	"testing"

	"github.com/issue9/assert/v3"
	"golang.org/x/text/language"
)

func TestMergeLanguage(t *testing.T) {
	a := assert.New(t, false)
	log := make(chan string, 1000)

	src := &Language{
		ID:       language.SimplifiedChinese,
		Messages: []Message{{Key: "src"}},
	}
	dest := &Language{
		ID:       language.Afrikaans,
		Messages: []Message{{Key: "dest"}},
	}
	mergeLanguage(src, dest, log)
	a.Zero(len(log)).
		Length(src.Messages, 1).Equal(src.Messages[0].Key, "src").
		Length(dest.Messages, 1).Equal(dest.Messages[0].Key, "dest")

	src = &Language{
		ID:       language.SimplifiedChinese,
		Messages: []Message{{Key: "src"}},
	}
	dest = &Language{
		ID:       language.SimplifiedChinese,
		Messages: []Message{{Key: "dest"}},
	}
	mergeLanguage(src, dest, log)
	a.Equal(len(log), 1).
		Length(src.Messages, 1).Equal(src.Messages[0].Key, "src").
		Length(dest.Messages, 1).Equal(dest.Messages[0].Key, "src")

	<-log
	src = &Language{
		ID:       language.SimplifiedChinese,
		Messages: []Message{{Key: "src"}, {Key: "g"}},
	}
	dest = &Language{
		ID:       language.SimplifiedChinese,
		Messages: []Message{{Key: "dest"}, {Key: "g"}},
	}
	mergeLanguage(src, dest, log)
	a.Equal(len(log), 1).
		Length(src.Messages, 2).
		Length(dest.Messages, 2).Equal(dest.Messages[0].Key, "g").Equal(dest.Messages[1].Key, "src")
}

func TestMerge(t *testing.T) {
	a := assert.New(t, false)
	log := make(chan string, 1000)

	src := &Messages{
		Languages: []*Language{{ID: language.SimplifiedChinese}},
	}
	dest := &Messages{
		Languages: []*Language{{ID: language.Afrikaans}},
	}
	dest.Merge(src, log)
	a.Equal(len(log), 1).
		Length(src.Languages, 1).Equal(src.Languages[0].ID, language.SimplifiedChinese).
		Length(dest.Languages, 1).Equal(dest.Languages[0].ID, language.SimplifiedChinese)

	<-log
	src = &Messages{
		Languages: []*Language{{
			ID:       language.SimplifiedChinese,
			Messages: []Message{{Key: "src"}},
		}},
	}
	dest = &Messages{
		Languages: []*Language{{
			ID:       language.SimplifiedChinese,
			Messages: []Message{{Key: "dest"}},
		}},
	}
	dest.Merge(src, log)
	a.Equal(len(log), 1).
		Length(src.Languages, 1).Equal(src.Languages[0].ID, language.SimplifiedChinese).
		Length(dest.Languages, 1).Equal(dest.Languages[0].ID, language.SimplifiedChinese).
		Length(dest.Languages[0].Messages, 1).Equal(dest.Languages[0].Messages[0].Key, "src")
}
