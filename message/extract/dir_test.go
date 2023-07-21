// SPDX-License-Identifier: MIT

package extract

import (
	"testing"

	"github.com/issue9/assert/v3"
)

func TestGetDirs(t *testing.T) {
	a := assert.New(t, false)

	dirs, err := getDir("./", false, false)
	a.NotError(err).Length(dirs, 1)

	dirs, err = getDir("./", false, true)
	a.NotError(err).Length(dirs, 1)

	dirs, err = getDir("./", true, false)
	a.NotError(err).Length(dirs, 3, "%+v", dirs)

	dirs, err = getDir("./", true, true)
	a.NotError(err).Length(dirs, 1)

	dirs, err = getDir("./testdata", true, true)
	a.NotError(err).Length(dirs, 2)
}
