// SPDX-License-Identifier: MIT

//go:build !windows && !js
// +build !windows,!js

package syslocale

func getLocaleName() (string, error) { return getEnvLang(), nil }
