// SPDX-License-Identifier: MIT

package message

import (
	"golang.org/x/text/feature/plural"
	"golang.org/x/text/message/catalog"
)

func ex(cases []*Case) []any {
	data := make([]any, 0, len(cases)*2)
	for _, c := range cases {
		data = append(data, c.Case, c.Value)
	}
	return data
}

// Catalog 将当前对象附加在 [catalog.Catalog] 上
func (m *Messages) Catalog(b *catalog.Builder) (err error) {
	for _, lang := range m.Languages {
		for _, msg := range lang.Messages {
			switch {
			case msg.Message.Vars != nil:
				vars := msg.Message.Vars
				msgs := make([]catalog.Message, 0, len(vars))
				for _, v := range vars {
					mm := catalog.Var(v.Name, plural.Selectf(v.Arg, v.Format, ex(v.Cases)...))
					msgs = append(msgs, mm)
				}
				msgs = append(msgs, catalog.String(msg.Message.Msg))
				err = b.Set(lang.ID, msg.Key, msgs...)
			case msg.Message.Select != nil:
				s := msg.Message.Select
				err = b.Set(lang.ID, msg.Key, plural.Selectf(s.Arg, s.Format, ex(s.Cases)...))
			case msg.Message.Msg != "":
				err = b.SetString(lang.ID, msg.Key, msg.Message.Msg)
			}

			if err != nil {
				return err
			}
		}
	}

	return nil
}

