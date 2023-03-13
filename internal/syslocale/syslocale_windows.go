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

func getOSLocaleName() string {
	k32, err := syscall.LoadDLL("kernel32.dll")
	if err != nil {
		log.Println(err)
		return ""
	}
	defer k32.Release()

	f, err := k32.FindProc(getLocaleFuncName)
	if err != nil {
		log.Println(err)
		return ""
	}

	buf := make([]uint16, maxLen)
	r1, _, err := f.Call(uintptr(unsafe.Pointer(&buf[0])), uintptr(maxLen))
	if uint32(r1) == 0 {
		log.Println(err)
		return ""
	}

	return syscall.UTF16ToString(buf)
}
