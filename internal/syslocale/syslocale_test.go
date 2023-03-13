// SPDX-License-Identifier: MIT

package syslocale

import (
	"testing"

	"github.com/issue9/assert/v3"
)

func TestGetOSLocaleName(t *testing.T) {
	a := assert.New(t, false)

	name := getOSLocaleName()
	a.True(len(name) > 0)
}
