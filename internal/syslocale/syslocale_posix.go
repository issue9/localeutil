// SPDX-License-Identifier: MIT

//go:build !windows && !js && !darwin
// +build !windows,!js,!darwin

package syslocale

func getLocaleName() (string, error) { return getEnvLang(), nil }
