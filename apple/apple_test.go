// SPDX-FileCopyrightText: 2020-2024 caixw
//
// SPDX-License-Identifier: MIT

//go:build darwin || ios

package apple

import (
	"testing"

	"github.com/issue9/assert/v4"
)

func TestAppLocale(t *testing.T) {
	a := assert.New(t, false)

	app := "com.example.test"

	a.NotError(SetAppLocale(app, "zh-CN", "zh-TW"))
	id := AppLocale(app)
	a.Equal(id, `zh-CN`)
}
