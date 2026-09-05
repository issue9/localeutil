// SPDX-FileCopyrightText: 2020-2026 caixw
//
// SPDX-License-Identifier: MIT

package syslocale

import (
	"testing"

	"github.com/issue9/assert/v5"
)

func TestGet(t *testing.T) {
	a := assert.New(t, false)

	name := Get()
	a.True(len(name) > 0)
}
