// SPDX-FileCopyrightText: 2020-2024 caixw
//
// SPDX-License-Identifier: MIT

package syslocale

import (
	"testing"

	"github.com/issue9/assert/v4"
)

func TestGet(t *testing.T) {
	a := assert.New(t, false)

	name := Get()
	a.True(len(name) > 0)
}
