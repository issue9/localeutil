// SPDX-License-Identifier: MIT

//go:build darwin || ios

package defaults

import (
	"testing"

	"github.com/issue9/assert/v3"
)

func TestReadWrite(t *testing.T) {
	a := assert.New(t, false)

	const key = "key"
	const domain = "com.example.test"

	a.NotError(Write(domain, key, "string", "123"))

	v := ReadDomains(key, domain)
	a.Equal(v, "123")

	v = ReadDomains("not-exists", domain)
	a.Equal(v, "")

	v = ReadDomains(key, "not-exists")
	a.Equal(v, "")
}
