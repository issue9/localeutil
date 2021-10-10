// SPDX-License-Identifier: MIT

//go:build !windows
// +build !windows

package syslocale

func getLocaleName() (string, error) { return getEnvLang(), nil }
