// SPDX-License-Identifier: MIT

//go:build !windows && !js && !darwin
// +build !windows,!js,!darwin

package syslocale

func getOSLocaleName() (string, error) { return "", nil }
