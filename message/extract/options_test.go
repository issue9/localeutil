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
	a.NotError(err).Length(dirs, 4, "%+v", dirs)

	dirs, err = getDir("./", true, true)
	a.NotError(err).Length(dirs, 1)

	dirs, err = getDir("./testdata", true, true)
	a.NotError(err).Length(dirs, 3)
}

func TestSplit(t *testing.T) {
	a := assert.New(t, false)

	fns := split("github.com/issue9/localeutil.Phrase", "github.com/issue9/localeutil.Error", "github.com/issue9/localeutil.Struct.Printf")
	a.Equal(fns, []fn{
		{pkgName: "github.com/issue9/localeutil", name: "Phrase"},
		{pkgName: "github.com/issue9/localeutil", name: "Error"},
		{pkgName: "github.com/issue9/localeutil", name: "Printf", typeName: "Struct"},
	})

	a.PanicString(func() {
		split("github.com/issue9")
	}, "github.com/issue9 格式无效")
}
