// SPDX-License-Identifier: MIT

package syslocale

import (
	"syscall"
	"unsafe"
)

const (
	getLocaleFuncName = "GetUserDefaultLocaleName"
	maxLen            = 85 // GetUserDefaultLocaleName 第二个参数
)

func getLocaleName() (string, error) {
	if name := getEnvLang(); len(name) > 0 {
		return name, nil
	}

	k32, err := syscall.LoadDLL("kernel32.dll")
	if err != nil {
		return "", err
	}
	defer k32.Release()

	f, err := k32.FindProc(getLocaleFuncName)
	if err != nil {
		return "", err
	}

	buf := make([]uint16, maxLen)
	r1, _, err := f.Call(uintptr(unsafe.Pointer(&buf[0])), uintptr(maxLen))
	if uint32(r1) == 0 {
		return "", err
	}

	return syscall.UTF16ToString(buf), nil
}
