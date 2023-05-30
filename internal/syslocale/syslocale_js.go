// SPDX-License-Identifier: MIT

package syslocale

import "syscall/js"

func getNavigator() js.Value { return js.Global().Get("navigator") }

func getOSLocaleName() string {
	nav := getNavigator()
	if nav.IsUndefined() {
		return ""
	}

	lang := nav.Get("language")
	if lang.IsUndefined() {
		return ""
	}

	return lang.String()
}
