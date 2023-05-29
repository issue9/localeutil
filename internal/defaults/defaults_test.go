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

	v, err := ReadDomains(key, domain)
	a.NotError(err).Equal(v, "123")
}
