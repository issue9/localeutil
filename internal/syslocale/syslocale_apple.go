// SPDX-License-Identifier: MIT

//go:build darwin || ios

package syslocale

import "github.com/issue9/localeutil/internal/defaults"

func getOSLocaleName() string {
	return defaults.ReadDomains(
		"AppleLocale",
		"~/Library/Preferences/.GlobalPreferences",
		"/Library/Preferences/.GlobalPreferences",
		"-g",
	)
}
