// SPDX-License-Identifier: MIT

package syslocale

import (
	"testing"

	"github.com/issue9/assert/v3"
)

func TestGet(t *testing.T) {
	a := assert.New(t, false)

	name := Get()
	a.True(len(name) > 0)
}
