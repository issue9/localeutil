// SPDX-License-Identifier: MIT

package windows

import (
	"testing"

	"github.com/issue9/assert/v3"
)

func TestGetLCID(t *testing.T) {
	a := assert.New(t, false)

	// 测试数据
	// https://learn.microsoft.com/en-us/openspecs/office_standards/ms-oe376/6c085406-a698-4e12-9d4d-c3b0ee3dbc4a

	lcid, err := GetLCID("zh-CN")
	a.NotError(err).Equal(lcid, 2052)

	lcid, err = GetLCID("zh-hans")
	a.NotError(err).Equal(lcid, 2052)

	lcid, err = GetLCID("cmn-hans")
	a.NotError(err).Equal(lcid, 4096)
}
