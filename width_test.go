// SPDX-License-Identifier: MIT

package localeutil

import (
	"testing"

	"github.com/issue9/assert/v3"
)

func TestWidth(t *testing.T) {
	a := assert.New(t, false)

	a.Equal(4, Width("汉字"))
	a.Equal(3, Width("3,a"))
	a.Equal(10, Width("汉字３，Ａ"))
}
