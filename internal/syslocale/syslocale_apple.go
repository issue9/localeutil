// SPDX-License-Identifier: MIT

//go:build darwin || ios

package syslocale

import (
	"log"

	"github.com/issue9/localeutil/internal/defaults"
)

func getOSLocaleName() string {
	val, err := defaults.ReadDomains(
		"AppleLocale",
		"~/Library/Preferences/.GlobalPreferences",
		"/Library/Preferences/.GlobalPreferences",
		"-g",
	)
	if err != nil {
		log.Println(err)
	}
	return val
}
