// SPDX-License-Identifier: MIT

package apple

import (
	"testing"

	"github.com/issue9/assert/v3"
)

func TestAppLocale(t *testing.T) {
	a := assert.New(t, false)

	app := "com.example.test"

	a.NotError(SetAppLocale(app, "zh-CN", "zh-TW"))
	id, err := AppLocale(app)
	a.NotError(err).Equal(id, `"zh-CN"`)
}
