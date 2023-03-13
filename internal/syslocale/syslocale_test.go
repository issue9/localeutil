// SPDX-License-Identifier: MIT

package syslocale

import (
	"testing"

	"github.com/issue9/assert/v3"
)

func TestGetLocaleName(t *testing.T) {
	a := assert.New(t, false)

	name, err := getLocaleName()
	a.NotError(err).True(len(name) > 0)
}
