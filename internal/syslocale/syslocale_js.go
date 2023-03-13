// SPDX-License-Identifier: MIT

package syslocale

import (
	"errors"
	"syscall/js"
)

func getNavigator() js.Value { return js.Global().Get("navigator") }

func getOSLocaleName() string {
	nav := getNavigator()
	if nav.IsUndefined() {
		log.Println("未定义 window.navigator")
		return ""
	}

	lang := nav.Get("language")
	if lang.IsUndefined() {
		log.Println("未定义 window.navigator.language")
		return ""
	}

	return lang.String()
}
