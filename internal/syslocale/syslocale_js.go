// SPDX-License-Identifier: MIT

package syslocale

import (
	"errors"
	"syscall/js"
)

func getNavigator() js.Value {
	return js.Global().Get("navigator")
}

func getLocaleName() (string, error) {
	nav := getNavigator()
	if nav.IsUndefined() {
		return "", errors.New("未定义 window.navigator")
	}

	lang := nav.Get("language")
	if lang.IsUndefined() {
		return "", errors.New("未定义 window.navigator.language")
	}

	return lang.String(), nil
}
